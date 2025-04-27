package handlers

import (
	"banking_ledger/logger"
	"banking_ledger/models"
	"banking_ledger/services"
	"banking_ledger/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {

	var input models.CreateAccountRequest

	ctx := utils.GetContextFromGinContext(c)

	err := c.BindJSON(&input)
	if err != nil {
		errMsg := fmt.Sprintf("CreateAccount: Request body validation fail.Request body:%s.Error:%s", utils.ConvertStructToString(input), err.Error())
		logger.Log.Error(errMsg)
		apiError := utils.RenderApiError(ctx, http.StatusBadRequest, 2001, errMsg, "", nil)
		ProcessError(ctx, models.API_ERROR_NO_INTERVENTION_REQUIRED, errMsg, apiError)
		c.JSON(apiError.StatusCode, apiError.ApplicationError)
		return
	}

	user_id, exists := c.Get("user_id")
	if !exists {
		errMsg := "CreateAccount: UserId not found in context claims"
		logger.Log.Error(errMsg)
		apiError := utils.RenderApiError(ctx, http.StatusBadRequest, 2001, errMsg, "", nil)
		ProcessError(ctx, models.API_ERROR_NO_INTERVENTION_REQUIRED, errMsg, apiError)
		c.JSON(apiError.StatusCode, apiError.ApplicationError)
		return
	}

	userId, ok := user_id.(int)
	if !ok {
		errMsg := "CreateAccount: UserId found in context claims is not of correct type"
		logger.Log.Error(errMsg)
		apiError := utils.RenderApiError(ctx, http.StatusBadRequest, 2001, errMsg, "", nil)
		ProcessError(ctx, models.API_ERROR_NO_INTERVENTION_REQUIRED, errMsg, apiError)
		c.JSON(apiError.StatusCode, apiError.ApplicationError)
		return
	}

	apiError := services.CreateAccountForUser(ctx, userId, input)
	if apiError != nil {
		c.JSON(apiError.StatusCode, apiError.ApplicationError)
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Type: "success", Message: "Account created successfully"})
}
