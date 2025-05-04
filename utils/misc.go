package utils

import (
	"banking_ledger/logger"
	"banking_ledger/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetContextFromGinContext(c *gin.Context) (ctx context.Context) {

	ctx = c.Request.Context()

	return ctx

}

func ConvertStructToString(structInput interface{}) string {

	structInputByteArray, err := json.Marshal(structInput)
	if err != nil {
		errorMsg := fmt.Sprintf("Could not convert struct input to byte array.Struct value:%v.Error message:%s!", structInput, err.Error())
		logger.Log.Error(errorMsg)
		return ""
	}

	return string(structInputByteArray)

}

func CreateContextWithNewRequestId() (ctx context.Context) {

	requestId := uuid.New().String()
	ctx = context.WithValue(context.Background(), models.CONTEXT_REQUEST_ID_KEY, requestId)

	return ctx

}

func GetClaimFromContext[T any](c *gin.Context, key string) (T, error) {
	val, exists := c.Get(key)
	if !exists {
		return *new(T), fmt.Errorf("%s not found in context", key)
	}

	castVal, ok := val.(T)
	if !ok {
		return *new(T), fmt.Errorf("%s in context has wrong type", key)
	}

	return castVal, nil
}
