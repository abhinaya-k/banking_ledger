package misc

import (
	"banking_ledger/database"
	"banking_ledger/logger"
	"banking_ledger/models"
	"context"
	"fmt"
)

func ProcessError(ctx context.Context, priority int, errorMessage string, additionalInfo interface{}) {

	switch value := additionalInfo.(type) {

	case []byte:

		additionalInfoConverted, ok := additionalInfo.([]byte)

		if ok {

			appError := database.ErDb.ProcessErrorMessages(ctx, priority, errorMessage, string(additionalInfoConverted))

			if appError != nil {

				errorMsg := fmt.Sprintf("Failed to write ProcessError message to database.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
				logger.Log.Error(errorMsg)
			}

		} else {

			errorMsg := fmt.Sprintf("Could not write ProcessError message to database as byte array could not be converted to string.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
			logger.Log.Error(errorMsg)

		}

	case *models.ApplicationError:

		additionalInfoConverted, ok := additionalInfo.(*models.ApplicationError)

		if ok {

			additionalInfoToWriteToDb := fmt.Sprintf("ApplicationErrorMessage Details-ErrorCode:%d,ErrorMessage:%s", additionalInfoConverted.Message.ErrorCode, additionalInfoConverted.Message.ErrorMessage)

			appError := database.ErDb.ProcessErrorMessages(ctx, priority, errorMessage, additionalInfoToWriteToDb)

			if appError != nil {

				errorMsg := fmt.Sprintf("Failed to write ProcessError message to database.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
				logger.Log.Error(errorMsg)
			}

		} else {

			errorMsg := fmt.Sprintf("Could not write ProcessError message to database as *models.ApplicationError message could not be converted.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
			logger.Log.Error(errorMsg)

		}

	case *models.ApiError:

		additionalInfoConverted, ok := additionalInfo.(*models.ApiError)

		if ok {

			additionalInfoToWriteToDb := fmt.Sprintf("ApiErrorMessage Details-StatusCode:%d,ErrorCode:%d,ErrorMessage:%s", additionalInfoConverted.StatusCode, additionalInfoConverted.ApplicationError.Message.ErrorCode, additionalInfoConverted.ApplicationError.Message.ErrorMessage)

			appError := database.ErDb.ProcessErrorMessages(ctx, priority, errorMessage, additionalInfoToWriteToDb)

			if appError != nil {

				errorMsg := fmt.Sprintf("Failed to write ProcessError message to database.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
				logger.Log.Error(errorMsg)
			}

		} else {

			errorMsg := fmt.Sprintf("Could not write ProcessError message to database as *models.ApplicationError message could not be converted.Message priority:%d,Message error:%s,Message data:%v", priority, errorMessage, additionalInfo)
			logger.Log.Error(errorMsg)

		}

	case nil:

		appError := database.ErDb.ProcessErrorMessages(ctx, priority, errorMessage, "")

		if appError != nil {

			errorMsg := fmt.Sprintf("Failed to write ProcessError message to database for nil message.Message priority:%d,Message error:%s", priority, errorMessage)
			logger.Log.Error(errorMsg)
		}

	default:
		errorMsg := fmt.Sprintf("Could not process ProcessError message as type of data interface is not supported.Data interface type:%v,Message priority:%d,Message error:%s,Message data:%v", value, priority, errorMessage, additionalInfo)
		logger.Log.Error(errorMsg)

	}

}
