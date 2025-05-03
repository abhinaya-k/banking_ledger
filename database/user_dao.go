package database

import (
	"banking_ledger/logger"
	"banking_ledger/models"
	"banking_ledger/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type userDb struct{}

type userDbInterface interface {
	GetUserByEmail(ctx context.Context, email string) (exists bool, user models.User, appError *models.ApplicationError)
	CreateUser(ctx context.Context, userDetails models.User) *models.ApplicationError
	GetUserByUserId(ctx context.Context, userId int) (exists bool, user models.User, appError *models.ApplicationError)
}

var UserDb userDbInterface

func init() {
	UserDb = &userDb{}
}

func (u *userDb) GetUserByEmail(ctx context.Context, email string) (exists bool, user models.User, appError *models.ApplicationError) {

	sqlStatement := `select u."user_id", u."email", u."password_hash", u."first_name", u."last_name" from users u where  u."email" = $1`

	err := dbPool.QueryRow(ctx, sqlStatement, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName)
	if err != nil {

		if err == pgx.ErrNoRows {
			return false, user, nil
		}

		errMsg := fmt.Sprintf("GetUserByEmail: Could not get user details from Database. Error:%s!", err.Error())
		displayMsg := fmt.Sprintf("Could not get user details for emailId: %s!", email)
		logger.Log.Error(errMsg)
		appError = utils.RenderAppError(ctx, 1001, errMsg, displayMsg, nil)
		return false, user, appError
	}

	return true, user, nil
}

func (u *userDb) CreateUser(ctx context.Context, userDetails models.User) *models.ApplicationError {

	var userId int

	sqlStatement := `INSERT INTO users ( "email", "password_hash", "first_name","last_name") VALUES ($1, $2, $3, $4) RETURNING user_id;`

	err := dbPool.QueryRow(ctx, sqlStatement, userDetails.Email, userDetails.PasswordHash, userDetails.FirstName, userDetails.LastName).Scan(&userId)
	if err != nil {
		errMsg := fmt.Sprintf("CreateUser: Couldn't insert user details. Error:%s!", err.Error())
		displayMsg := fmt.Sprintf("Could not save user details for emailId: %s", userDetails.Email)
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 1001, errMsg, displayMsg, nil)
		return appError
	}

	return nil
}

func (u *userDb) GetUserByUserId(ctx context.Context, userId int) (exists bool, user models.User, appError *models.ApplicationError) {

	sqlStatement := `select u."user_id", u."email", u."password_hash", u."first_name", u."last_name" from users u where  u."user_id" = $1`

	err := dbPool.QueryRow(ctx, sqlStatement, userId).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName)
	if err != nil {

		if err == pgx.ErrNoRows {
			return false, user, nil
		}

		errMsg := fmt.Sprintf("GetUserByUserId: Could not get user details from Database. Error:%s!", err.Error())
		displayMsg := fmt.Sprintf("Could not get user details for userId: %d!", userId)
		logger.Log.Error(errMsg)
		appError = utils.RenderAppError(ctx, 1001, errMsg, displayMsg, nil)
		return false, user, appError
	}

	return true, user, nil
}
