package services

import (
	"banking_ledger/clients"
	"banking_ledger/config"
	"banking_ledger/database"
	"banking_ledger/logger"
	"banking_ledger/misc"
	"banking_ledger/models"
	"banking_ledger/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateAccountForUser(ctx context.Context, userId int, req models.CreateAccountRequest) *models.ApiError {

	tx, err := database.AccDb.BeginTx(ctx)
	if err != nil {
		errMsg := "CreateAccountForUser: Could not begin transaction!"
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 1001, errMsg, "", nil)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	transactionErrMsg := "Transaction failed"
	txCommitted := false

	defer func() {

		if !txCommitted {

			tx.Rollback(ctx)

			transactionToLog := models.TransactionCollection{
				UserId:            userId,
				Amount:            req.Balance,
				TransactionType:   "deposit",
				TransactionStatus: "failed",
				TransactionMsg:    transactionErrMsg,
				RequestId:         uuid.New(),
				TransactionTime:   time.Now().Unix(),
			}

			txCollection := database.GetCollection("transactions")

			_, err = txCollection.InsertOne(ctx, transactionToLog)
			if err != nil {
				errMsg := fmt.Sprintf("CreateAccountForUser: Failed to insert transaction into MongoDB! Error: %s", err.Error())
				logger.Log.Error(errMsg)
				appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
				misc.ProcessError(ctx, models.KAFKA_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
			}
		}

	}()

	exists, _, appError := database.AccDb.GetAccountByUserId(ctx, tx, userId)
	if appError != nil {
		transactionErrMsg = "Internal Error"
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, "CreateAccountForUser-> Failed to check if account exists", appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if exists {
		errMsg := "Account already exists for this user"
		logger.Log.Error(errMsg)
		return utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	balanceInPaise := int64(req.Balance * 100)

	appError = database.AccDb.CreateAccountForUser(ctx, tx, userId, balanceInPaise)
	if appError != nil {
		transactionErrMsg = "Internal Error"
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, "CreateAccountForUser-> Failed to create account for user", appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	transactionToLog := models.TransactionCollection{
		UserId:            userId,
		Amount:            req.Balance,
		TransactionType:   "deposit",
		TransactionStatus: "success",
		TransactionMsg:    "Account created successfully",
		RequestId:         uuid.New(),
		TransactionTime:   time.Now().Unix(),
	}

	txCollection := database.GetCollection("transactions")

	_, err = txCollection.InsertOne(ctx, transactionToLog)
	if err != nil {
		errMsg := fmt.Sprintf("CreateAccountForUser: Failed to insert transaction into MongoDB! Error: %s", err.Error())
		logger.Log.Error(errMsg)
		transactionErrMsg = "Internal Error"
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if err := tx.Commit(ctx); err != nil {
		errMsg := "CreateAccountForUser: Failed to commit transaction!"
		logger.Log.Error(errMsg)
		transactionErrMsg = "Internal Error"
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	txCommitted = true

	return nil

}

func FundTransaction(ctx context.Context, userId int, req models.FundTransactionRequest) *models.ApiError {

	tx, err := database.AccDb.BeginTx(ctx)
	if err != nil {
		errMsg := "FundTransaction: Could not begin transaction!"
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 1001, errMsg, "", nil)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	defer tx.Rollback(ctx)

	exists, _, appError := database.AccDb.GetAccountByUserId(ctx, tx, userId)
	if appError != nil {
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, "Failed to check if account exists", appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if !exists {
		errMsg := "Account does not exists for this user!"
		logger.Log.Error(errMsg)
		return utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	kafkaMsg := models.TransactionRequestKafka{
		UserId:          userId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		RequestId:       uuid.New(),
		TransactionTime: time.Now().Unix(),
	}

	appError = clients.SendMessageToKafkaTopic(ctx, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, kafkaMsg, strconv.Itoa(userId))
	if appError != nil {
		errMsg := fmt.Sprintf("FundTransaction: Failed to send message to Kafka topic! Error: %s", appError.Message.ErrorMessage)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if err := tx.Commit(ctx); err != nil {
		errMsg := "FundTransaction: Failed to commit transaction!"
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	return nil

}

func ProcessTransaction(ctx context.Context, transaction models.TransactionRequestKafka) *models.ApplicationError {

	tx, err := database.AccDb.BeginTx(ctx)
	if err != nil {
		errMsg := "ProcessTransaction: Could not begin transaction!"
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 1001, errMsg, "", nil)
		misc.ProcessError(ctx, models.KAFKA_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		return appError
	}

	transactionErrMsg := "Transaction failed"
	txCommitted := false

	defer func() {

		if !txCommitted {

			tx.Rollback(ctx)

			transactionToLog := models.TransactionCollection{
				UserId:            transaction.UserId,
				Amount:            transaction.Amount,
				TransactionType:   transaction.TransactionType,
				TransactionStatus: "failed",
				TransactionMsg:    transactionErrMsg,
				RequestId:         transaction.RequestId,
				TransactionTime:   transaction.TransactionTime,
			}

			txCollection := database.GetCollection("transactions")

			_, err = txCollection.InsertOne(ctx, transactionToLog)
			if err != nil {
				errMsg := fmt.Sprintf("ProcessTransaction: Failed to insert transaction into MongoDB! Error: %s", err.Error())
				logger.Log.Error(errMsg)
				appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
				misc.ProcessError(ctx, models.KAFKA_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
			}
		}

	}()

	exists, balance, appError := database.AccDb.GetBalanceForUserId(ctx, tx, transaction.UserId)
	if appError != nil {
		errMsg := fmt.Sprintf("ProcessTransaction: Failed to get balance for user %d", transaction.UserId)
		logger.Log.Error(errMsg)
		transactionErrMsg = "Failed to get balance for user"
		misc.ProcessError(ctx, models.KAFKA_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		return appError
	}

	if !exists {
		errMsg := fmt.Sprintf("ProcessTransaction: Account details not found for user! UserId: %d", transaction.UserId)
		logger.Log.Error(errMsg)
		transactionErrMsg = "Failed to get balance for user"
		appError := utils.RenderAppError(ctx, 1001, errMsg, "", nil)
		return appError
	}

	if transaction.TransactionType == "withdraw" && balance < int64(transaction.Amount*100) {
		errMsg := fmt.Sprintf("ProcessTransaction: Insufficient balance for user! UserId: %d", transaction.UserId)
		logger.Log.Error(errMsg)
		transactionErrMsg = "Insufficient balance for user!"
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return appError
	}

	var newBalance int64
	switch transaction.TransactionType {

	case "deposit":
		newBalance = balance + int64(transaction.Amount*100)
	case "withdraw":
		newBalance = balance - int64(transaction.Amount*100)
	default:
		errMsg := fmt.Sprintf("ProcessTransaction: Invalid Transaction Type! TransactionType: %s", transaction.TransactionType)
		logger.Log.Error(errMsg)
		transactionErrMsg = "Invalid Request!"
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return appError
	}

	appError = database.AccDb.UpdateBalanceForUserId(ctx, tx, transaction.UserId, newBalance)
	if appError != nil {
		transactionErrMsg = "Internal Error!"
		misc.ProcessError(ctx, models.KAFKA_ERROR_REQUIRE_INTERVENTION, "Failed to update balance for user", appError)
		return appError
	}

	transactionToLog := models.TransactionCollection{
		UserId:            transaction.UserId,
		Amount:            transaction.Amount,
		TransactionType:   transaction.TransactionType,
		TransactionStatus: "success",
		TransactionMsg:    "Transaction completed successfully",
		RequestId:         transaction.RequestId,
		TransactionTime:   transaction.TransactionTime,
	}

	txCollection := database.GetCollection("transactions")

	_, err = txCollection.InsertOne(ctx, transactionToLog)
	if err != nil {
		errMsg := fmt.Sprintf("ProcessTransaction: Failed to insert transaction into MongoDB! Error: %s", err.Error())
		logger.Log.Error(errMsg)
		transactionErrMsg = "Internal Error!"
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return appError
	}

	if err := tx.Commit(ctx); err != nil {
		errMsg := "ProcessTransaction: Failed to commit transaction!"
		logger.Log.Error(errMsg)
		transactionErrMsg = "Internal Error!"
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return appError
	}

	txCommitted = true

	return nil

}

func GetTransactionHistory(ctx context.Context, userId int, req models.GetTransactionHistoryRequest, role string) (*models.GetTransactionHistoryResponse, *models.ApiError) {

	tx, err := database.AccDb.BeginTx(ctx)
	if err != nil {
		errMsg := "ProcessTransaction: Could not begin transaction!"
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 1001, errMsg, "", nil)
		misc.ProcessError(ctx, models.KAFKA_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		return nil, utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	defer tx.Rollback(ctx)

	exists, _, appError := database.AccDb.GetAccountByUserId(ctx, tx, userId)
	if appError != nil {
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, "Failed to check if account exists", appError)
		return nil, utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if !exists {
		errMsg := "Account does not exists for this user!"
		logger.Log.Error(errMsg)
		return nil, utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	userExists, user, appError := database.UserDb.GetUserByUserId(ctx, userId)
	if appError != nil {
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, "Failed to check if user exists", appError)
		return nil, utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if !userExists {
		errMsg := fmt.Sprintf("User does not exists UserId: %d!", userId)
		logger.Log.Error(errMsg)
		return nil, utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	txCollection := database.GetCollection("transactions")

	filter := bson.M{}

	if role != "admin" {
		filter["userId"] = userId
	}

	var deposit string = "deposit"
	var withdraw string = "withdraw"

	if req.Filters.TransactionType != nil && req.Filters.TransactionType != &deposit && req.Filters.TransactionType != &withdraw {
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

	fields := []zapcore.Field{
		zap.Any("filter", filter),
		zap.Any("findOptions ", findOptions),
	}

	logger.Log.Info("GetTransactionHistory: MongoDB query", fields...)

	cursor, err := txCollection.Find(ctx, filter, findOptions)
	if err != nil {
		errMsg := fmt.Sprintf("GetTransactionHistory: Failed to find transactions in MongoDB! Error: %s", err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 1111, errMsg, "", nil)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		apiError := utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
		return nil, apiError
	}
	defer cursor.Close(ctx)

	transactions := []models.TransactionHistory{}
	for cursor.Next(ctx) {
		var transaction models.TransactionCollection
		if err := cursor.Decode(&transaction); err != nil {
			errMsg := fmt.Sprintf("GetTransactionHistory: Failed to decode transaction! Error: %s", err.Error())
			logger.Log.Error(errMsg)
			appError := utils.RenderAppError(ctx, 1111, errMsg, "", nil)
			misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
			apiError := utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
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
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 1111, errMsg, "", nil)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, appError)
		apiError := utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
		return nil, apiError
	}

	if err := tx.Commit(ctx); err != nil {
		errMsg := "ProcessTransaction: Failed to commit transaction!"
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 1001, errMsg, errMsg, nil)
		return nil, utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	var apiResponse models.GetTransactionHistoryResponse

	apiResponse.TransactionHistory = transactions
	apiResponse.Pagination = *req.Pagination

	return &apiResponse, nil

}
