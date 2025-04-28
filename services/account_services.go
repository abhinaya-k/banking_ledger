package services

import (
	clients "banking_ledger/client"
	"banking_ledger/config"
	"banking_ledger/database"
	"banking_ledger/models"
	"banking_ledger/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateAccountForUser(ctx context.Context, userId int, req models.CreateAccountRequest) *models.ApiError {

	exists, _, appError := database.AccDb.GetAccountByUserId(ctx, userId)
	if appError != nil {
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if exists {
		errMsg := "Account already exists for this user"
		return utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	balanceInPaise := int(req.Balance * 100)

	appError = database.AccDb.CreateAccountForUser(ctx, userId, balanceInPaise)
	if appError != nil {
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	return nil

}

func FundTransaction(ctx context.Context, userId int, req models.FundTransactionRequest) *models.ApiError {

	if req.TransactionType != "deposit" || req.TransactionType != "withdraw" {
		errMsg := "Incorrect RequestBody! TransactionType should be either deposit or withdraw"
		return utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	exists, _, appError := database.AccDb.GetAccountByUserId(ctx, userId)
	if appError != nil {
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if !exists {
		errMsg := "Account does not exists for this user!"
		return utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	kafkaMsg := models.TransactionRequestKafka{
		UserId:          userId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		RequestId:       uuid.New(),
		TransactionTime: time.Now().Unix(),
	}

	appError = clients.SendMessageToKafkaTopic(ctx, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, kafkaMsg, string(userId))

	return nil

}

func ProcessTransaction(ctx context.Context, transaction models.TransactionRequestKafka) *models.ApplicationError {

	tx, err := database.AccDb.BeginTx(ctx)
	if err != nil {
		errMsg := "ProcessTransaction: Could not begin transaction!"
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return appError
	}

	defer tx.Rollback(ctx)

	exists, balance, appError := database.AccDb.GetBalanceForUserId(ctx, tx, transaction.UserId)
	if appError != nil {
		return appError
	}

	if !exists {
		errMsg := fmt.Sprintf("ProcessTransaction: Balance not for user! UserId: %d", transaction.UserId)
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return appError
	}

	if transaction.TransactionType == "withdraw" && balance < int(transaction.Amount*100) {
		errMsg := fmt.Sprintf("ProcessTransaction: Insufficient balance for user! UserId: %d", transaction.UserId)
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return appError
	}

	var newBalance int
	switch transaction.TransactionType {

	case "deposit":
		newBalance = balance + int(transaction.Amount*100)
	case "withdraw":
		newBalance = balance - int(transaction.Amount*100)
	default:
		errMsg := fmt.Sprintf("ProcessTransaction: Invalid Transaction Type! TransactionType: %s", transaction.TransactionType)
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return appError
	}

	appError = database.AccDb.UpdateBalanceForUserId(ctx, tx, transaction.UserId, newBalance)
	if appError != nil {
		return appError
	}

	txCollection := database.GetCollection("transactions")

	_, err = txCollection.InsertOne(ctx, transaction)
	if err != nil {
		errMsg := fmt.Sprintf("ProcessTransaction: Failed to insert transaction into MongoDB! Error: %s", err.Error())
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return appError
	}

	if err := tx.Commit(ctx); err != nil {
		errMsg := "ProcessTransaction: Failed to commit transaction!"
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return appError
	}

	return nil

}
