package services

import (
	"banking_ledger/database"
	"banking_ledger/logger"
	"banking_ledger/misc"
	"banking_ledger/models"
	"banking_ledger/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, req models.RegisterUserReqBody) *models.ApiError {
	// Validate the input
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		errMsg := "email and password are required"
		logger.Log.Error(errMsg)
		return utils.RenderApiError(ctx, http.StatusBadRequest, 5101, errMsg, "", nil)
	}

	// Check if user already exists
	exists, _, appError := database.UserDb.GetUserByEmail(ctx, req.Email)
	if appError != nil {
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, "Failed to check if user exists", appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if exists {
		errMsg := "user already exists with this email"
		logger.Log.Error(errMsg)
		return utils.RenderApiError(ctx, http.StatusBadRequest, 5102, errMsg, "", nil)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to generate password hash! Error: %s", err.Error())
		logger.Log.Error(errMsg)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, nil)
		return utils.RenderApiError(ctx, http.StatusInternalServerError, 5103, errMsg, "", nil)
	}

	// Create user model
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
	}

	// Save user to the database
	appError = database.UserDb.CreateUser(ctx, *user)
	if appError != nil {
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, "Failed to create user", appError)
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	return nil
}

func UserLogin(ctx context.Context, req models.LoginRequestBody) (response *models.LoginResponseBody, apiError *models.ApiError) {

	exists, user, appError := database.UserDb.GetUserByEmail(ctx, req.Email)
	if appError != nil {
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, "Failed to get user details", appError)
		return nil, utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if !exists {
		errMsg := fmt.Sprintf("user not found for email: %s", req.Email)
		displayMsg := "User not found! Please Register first"
		logger.Log.Error(errMsg)
		return nil, utils.RenderApiError(ctx, http.StatusBadRequest, 5104, errMsg, displayMsg, nil)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		errMsg := fmt.Sprintf("Error comparing password! Error: %s", err.Error())
		displayMsg := "Could not verify password"
		logger.Log.Error(errMsg)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, nil)
		return nil, utils.RenderApiError(ctx, http.StatusBadRequest, 5105, errMsg, displayMsg, nil)
	}

	fullName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	token, err := utils.GenerateJWTForUser(user.ID, user.Role, fullName)
	if err != nil {
		errMsg := fmt.Sprintf("Error generating token! Error: %s", err.Error())
		displayMsg := "Could not generate token"
		logger.Log.Error(errMsg)
		misc.ProcessError(ctx, models.API_ERROR_REQUIRE_INTERVENTION, errMsg, nil)
		return nil, utils.RenderApiError(ctx, http.StatusInternalServerError, 5106, errMsg, displayMsg, nil)
	}

	apiResponse := models.LoginResponseBody{
		Token: token,
	}

	return &apiResponse, nil

}
