package utils

import (
	"banking_ledger/logger"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
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
