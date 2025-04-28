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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetTransactionHistory(ctx context.Context, userId int, req models.GetTransactionHistoryRequest) (*models.GetTransactionHistoryResponse, *models.ApiError) {

	exists, _, appError := database.AccDb.GetAccountByUserId(ctx, userId)
	if appError != nil {
		return nil, utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if !exists {
		errMsg := "Account does not exists for this user!"
		return nil, utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	userExists, user, appError := database.UserDb.GetUserByUserId(ctx, userId)
	if appError != nil {
		return nil, utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if !userExists {
		errMsg := fmt.Sprintf("User does not exists UserId: %d!", userId)
		return nil, utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	txCollection := database.GetCollection("transactions")

	filter := bson.M{
		"userId": userId,
	}

	var deposit string = "deposit"
	var withdraw string = "withdraw"

	if req.Filters.TransactionType != nil && *&req.Filters.TransactionType != &deposit && *&req.Filters.TransactionType != &withdraw {
		filter["transactionType"] = *req.Filters.TransactionType
	}

	timeConditions := bson.M{}
	if req.Filters.StartTime != nil {
		timeConditions["$gte"] = *req.Filters.StartTime
	}
	if req.Filters.EndTime != nil {
		timeConditions["$lte"] = *req.Filters.EndTime
	}
	if len(timeConditions) > 0 {
		filter["transactionTime"] = timeConditions
	}

	if req.Pagination == nil {
		req.Pagination = &models.Pagination{
			Page:  1,
			Limit: 10,
		}
	}

	skip := (req.Pagination.Page - 1) * req.Pagination.Limit
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(req.Pagination.Limit).
		SetSort(bson.D{{Key: "transactionTime", Value: -1}})

	cursor, err := txCollection.Find(ctx, filter, findOptions)
	if err != nil {
		errMsg := fmt.Sprintf("GetTransactionHistory: Failed to find transactions in MongoDB! Error: %s", err.Error())
		apiError := utils.RenderApiError(ctx, http.StatusInternalServerError, 1001, errMsg, errMsg, nil)
		return nil, apiError
	}
	defer cursor.Close(ctx)

	// var transactions []models.TransactionRequestKafka
	var transactions []models.TransactionHistory
	for cursor.Next(ctx) {
		var transaction models.TransactionRequestKafka
		if err := cursor.Decode(&transaction); err != nil {
			errMsg := fmt.Sprintf("GetTransactionHistory: Failed to decode transaction! Error: %s", err.Error())
			apiError := utils.RenderApiError(ctx, http.StatusInternalServerError, 1001, errMsg, errMsg, nil)
			return nil, apiError
		}

		transactionHistory := models.TransactionHistory{
			UserId:          transaction.UserId,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Amount:          transaction.Amount,
			TransactionType: transaction.TransactionType,
			TransactionTime: transaction.TransactionTime,
		}

		transactions = append(transactions, transactionHistory)
	}

	if err := cursor.Err(); err != nil {
		errMsg := fmt.Sprintf("GetTransactionHistory: Cursor error! Error: %s", err.Error())
		apiError := utils.RenderApiError(ctx, http.StatusInternalServerError, 1001, errMsg, errMsg, nil)
		return nil, apiError
	}

	var apiResponse models.GetTransactionHistoryResponse

	apiResponse.TransactionHistory = transactions
	apiResponse.Pagination = *req.Pagination

	return &apiResponse, nil

}
