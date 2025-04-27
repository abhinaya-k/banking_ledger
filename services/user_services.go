package services

import (
	"banking_ledger/database"
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
		return utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	// Check if user already exists
	exists, _, appError := database.UserDb.GetUserByEmail(ctx, req.Email)
	if appError != nil {
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if exists {
		errMsg := "user already exists with this email"
		return utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, "", nil)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to generate password hash! Error: %s", err.Error())
		return utils.RenderApiError(ctx, http.StatusInternalServerError, 1001, errMsg, "", nil)
	}

	// Create user model
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	// Save user to the database
	appError = database.UserDb.CreateUser(ctx, *user)
	if appError != nil {
		return utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	return nil
}

func UserLogin(ctx context.Context, req models.LoginRequestBody) (response *models.LoginResponseBody, apiError *models.ApiError) {

	exists, user, appError := database.UserDb.GetUserByEmail(ctx, req.Email)
	if appError != nil {
		return nil, utils.RenderApiErrorFromAppError(http.StatusInternalServerError, appError)
	}

	if !exists {
		errMsg := "user not found"
		displayMsg := "User not found! Please Register first"
		return nil, utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, displayMsg, nil)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		errMsg := fmt.Sprintf("Error comparing password! Error: %s", err.Error())
		displayMsg := "Could not verify password"
		return nil, utils.RenderApiError(ctx, http.StatusBadRequest, 1001, errMsg, displayMsg, nil)
	}

	token, err := utils.GenerateJWTForUser(user.ID)
	if err != nil {
		errMsg := fmt.Sprintf("Error generating token! Error: %s", err.Error())
		displayMsg := "Could not generate token"
		return nil, utils.RenderApiError(ctx, http.StatusInternalServerError, 1001, errMsg, displayMsg, nil)
	}

	apiResponse := models.LoginResponseBody{
		Token: token,
	}

	return &apiResponse, nil

}
