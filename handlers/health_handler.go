package handlers

import (
	"banking_ledger/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHealth(c *gin.Context) {

	c.JSON(http.StatusOK, models.SuccessResponse{Type: "success", Message: "Service is healthy!!!"})

}

func NoRoute(c *gin.Context) {

	c.JSON(http.StatusNotFound, models.ApplicationError{Type: "error", Message: models.ApplicationErrorMessage{ErrorCode: 1000, ErrorMessage: "Route not found"}})

}
