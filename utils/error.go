package utils

import (
	"banking_ledger/models"
	"context"
)

func RenderApiErrorFromAppError(statusCode int, appError *models.ApplicationError) *models.ApiError {

	apiError := models.ApiError{
		StatusCode:       statusCode,
		ApplicationError: *appError,
	}

	return &apiError
}

func RenderApiError(ctx context.Context, statusCode int, errorCode int, errorMessage string, displayMessage string, additonalInfo interface{}) *models.ApiError {

	apiError := models.ApiError{
		StatusCode: statusCode,
		ApplicationError: models.ApplicationError{
			Type: "error",
			Message: models.ApplicationErrorMessage{
				ErrorCode:      errorCode,
				ErrorMessage:   errorMessage,
				DisplayMessage: displayMessage,
				AdditionalInfo: additonalInfo,
			},
		},
	}

	return &apiError
}

func RenderAppError(ctx context.Context, errorCode int, errorMessage string, displayMessage string, additonalInfo interface{}) *models.ApplicationError {

	appError := models.ApplicationError{
		Type: "error",
		Message: models.ApplicationErrorMessage{
			ErrorCode:      errorCode,
			ErrorMessage:   errorMessage,
			DisplayMessage: displayMessage,
			AdditionalInfo: additonalInfo,
		},
	}

	return &appError
}
