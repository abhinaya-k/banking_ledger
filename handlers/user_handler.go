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

func RegisterUser(c *gin.Context) {

	var input models.RegisterUserReqBody

	ctx := utils.GetContextFromGinContext(c)

	err := c.BindJSON(&input)
	if err != nil {
		errMsg := fmt.Sprintf("RegisterUser: Request body validation fail.Request body:%s.Error:%s", utils.ConvertStructToString(input), err.Error())
		logger.Log.Error(errMsg)
		apiError := utils.RenderApiError(ctx, http.StatusInternalServerError, 2001, errMsg, "", nil)
		ProcessError(ctx, models.API_ERROR_NO_INTERVENTION_REQUIRED, errMsg, apiError)
		c.JSON(apiError.StatusCode, apiError.ApplicationError)
		return
	}

	apiError := services.RegisterUser(ctx, input)
	if apiError != nil {
		c.JSON(apiError.StatusCode, apiError.ApplicationError)
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Type: "success", Message: "User Registered successfully"})
}
