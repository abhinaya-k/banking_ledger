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
