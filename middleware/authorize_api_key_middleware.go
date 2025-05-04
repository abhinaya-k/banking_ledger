package middleware

import (
	"banking_ledger/utils"
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type API_KEY_TYPE int

const (
	API_KEY API_KEY_TYPE = iota
)

var (
	apiKeys = map[API_KEY_TYPE]string{
		API_KEY: os.Getenv("API_KEY"),
	}
)

func AuthorizeApiKey(keyType API_KEY_TYPE) gin.HandlerFunc {
	return func(c *gin.Context) {

		xApiKey := c.GetHeader("x-api-key")

		requestId := uuid.New().String()
		type contextKey string
		const requestIDKey contextKey = "requestId"
		ctx := context.WithValue(c.Request.Context(), requestIDKey, requestId)

		if xApiKey == "" {
			apiError := utils.RenderApiError(ctx, http.StatusUnauthorized, 4001, "X-API-KEY missing in header!", "X-API-KEY missing in header!", nil)
			c.Abort()
			c.JSON(apiError.StatusCode, apiError.ApplicationError)
			return
		}

		expectedApiKey, ok := apiKeys[keyType]
		if !ok || xApiKey != expectedApiKey {
			apiError := utils.RenderApiError(ctx, http.StatusUnauthorized, 4002, "API key mismatch!", "API key mismatch!", nil)
			c.Abort()
			c.JSON(apiError.StatusCode, apiError.ApplicationError)
			return
		}

	}
}
