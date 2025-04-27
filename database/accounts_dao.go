package database

import (
	"banking_ledger/logger"
	"banking_ledger/models"
	"banking_ledger/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type accountDb struct{}

type accountDbInterface interface {
	GetAccountByUserId(ctx context.Context, userId int) (exists bool, account models.Account, appError *models.ApplicationError)
	CreateAccountForUser(ctx context.Context, userId int, balance int) *models.ApplicationError
}

var AccDb accountDbInterface

func init() {
	AccDb = &accountDb{}
}

func (a *accountDb) GetAccountByUserId(ctx context.Context, userId int) (exists bool, account models.Account, appError *models.ApplicationError) {

	sqlStatement := `select ac."account_id", ac."user_id", ac."balance" from accounts ac where ac."user_id" = $1`

	err := dbPool.QueryRow(ctx, sqlStatement, userId).Scan(&account.AccountID, &account.UserID, &account.Balance)
	if err != nil {

		if err == pgx.ErrNoRows {
			return false, account, nil
		}

		errMsg := fmt.Sprintf("GetAccountByUserId: Could not get account details from Database. Error:%s!", err.Error())
		displayMsg := "Could not get account details for the user!"
		logger.Log.Error(errMsg)
		appError = utils.RenderAppError(ctx, 1001, errMsg, displayMsg, nil)
		return false, account, appError
	}

	return true, account, nil
}

func (a *accountDb) CreateAccountForUser(ctx context.Context, userId int, balance int) *models.ApplicationError {

	var accountId int

	sqlStatement := `INSERT INTO accounts (user_id, balance) VALUES ($1, $2) RETURNING account_id;`

	err := dbPool.QueryRow(ctx, sqlStatement, userId, balance).Scan(&accountId)
	if err != nil {
		errMsg := fmt.Sprintf("CreateAccountForUser: Couldn't insert user account details. Error:%s!", err.Error())
		displayMsg := fmt.Sprintf("Could not create account for userId: %d", userId)
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 1001, errMsg, displayMsg, nil)
		return appError
	}

	return nil
}
