package services

import (
	"banking_ledger/database"
	"banking_ledger/models"
	"banking_ledger/utils"
	"context"
	"net/http"
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
		errMsg := "Incorrect RequestBody!"
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

	return nil

}
