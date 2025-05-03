package services

import (
	"banking_ledger/config"
	"banking_ledger/logger"
	"banking_ledger/misc"
	"banking_ledger/models"
	"banking_ledger/utils"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func KafkaConsumerProcessTransactions(msg *kafka.Message) (err error) {

	ctx := utils.CreateContextWithNewRequestId()

	var jsonData map[string]interface{}
	logger.Log.Debug("Incoming message is", zap.String("Message", string(msg.Value)))
	err = json.Unmarshal(msg.Value, &jsonData)
	if err != nil {
		errMsg := fmt.Sprintf("KafkaConsumerProcessTransactions:Could not unmarshal kafka message,Kafka topic:%s,Kafka message:%s!", config.TRANSACTION_PROCESSING_KAFKA_TOPIC, string(msg.Value))
		logger.Log.Error(errMsg)
		misc.SaveDroppedMessage(ctx, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, "JSON_UNMARSHAL_FAIL", msg.Value)
		return nil
	}

	userId, ok := jsonData["userId"].(float64)
	if !ok {
		errMsg := fmt.Sprintf("KafkaConsumerProcessTransactions:userId is not of correct type,Kafka topic:%s,Kafka message:%s!", config.TRANSACTION_PROCESSING_KAFKA_TOPIC, string(msg.Value))
		logger.Log.Error(errMsg)
		misc.SaveDroppedMessage(ctx, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, "INCORRECT_USERID", msg.Value)
		return nil
	}

	amount, ok := jsonData["amount"].(float64)
	if !ok {
		errMsg := fmt.Sprintf("KafkaConsumerProcessTransactions:amount is not of correct type,Kafka topic:%s,Kafka message:%s!", config.TRANSACTION_PROCESSING_KAFKA_TOPIC, string(msg.Value))
		logger.Log.Error(errMsg)
		misc.SaveDroppedMessage(ctx, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, "INCORRECT_AMOUNT", msg.Value)
		return nil
	}

	transactionType, ok := jsonData["transactionType"].(string)
	if !ok {
		errMsg := fmt.Sprintf("KafkaConsumerProcessTransactions:transactionType is not of correct type,Kafka topic:%s,Kafka message:%s!", config.TRANSACTION_PROCESSING_KAFKA_TOPIC, string(msg.Value))
		logger.Log.Error(errMsg)
		misc.SaveDroppedMessage(ctx, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, "INCORRECT_TRANSACTION_TYPE", msg.Value)
		return nil
	}

	requestIdStr, ok := jsonData["requestId"].(string)
	if !ok {
		errMsg := fmt.Sprintf("KafkaConsumerProcessTransactions:requestId is not of correct type,Kafka topic:%s,Kafka message:%s!", config.TRANSACTION_PROCESSING_KAFKA_TOPIC, string(msg.Value))
		logger.Log.Error(errMsg)
		misc.SaveDroppedMessage(ctx, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, "INCORRECT_REQUESTID", msg.Value)
		return nil
	}

	requestId, err := uuid.Parse(requestIdStr)
	if err != nil {
		errMsg := fmt.Sprintf("KafkaConsumerProcessTransactions:requestId is not of type uuid,Kafka topic:%s,Kafka message:%s!", config.TRANSACTION_PROCESSING_KAFKA_TOPIC, string(msg.Value))
		logger.Log.Error(errMsg)
		misc.SaveDroppedMessage(ctx, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, "INCORRECT_REQUEST_ID_TYPE", msg.Value)
		return nil
	}

	transactionTime, ok := jsonData["transactionTime"].(float64)
	if !ok {
		errMsg := fmt.Sprintf("KafkaConsumerProcessTransactions:transactionTime is not of correct type,Kafka topic:%s,Kafka message:%s!", config.TRANSACTION_PROCESSING_KAFKA_TOPIC, string(msg.Value))
		logger.Log.Error(errMsg)
		misc.SaveDroppedMessage(ctx, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, "INCORRECT_TRANSACTION_TIME", msg.Value)
		return nil
	}

	transactionRequest := models.TransactionRequestKafka{
		UserId:          int(userId),
		Amount:          amount,
		TransactionType: transactionType,
		RequestId:       requestId,
		TransactionTime: int64(transactionTime),
	}

	appError := ProcessTransaction(ctx, transactionRequest)
	if appError != nil {
		errMsg := fmt.Sprintf("KafkaConsumerProcessTransactions:Could not process transaction,Transaction request:%v,Error message:%s", transactionRequest, appError.Message.ErrorMessage)
		logger.Log.Error(errMsg)
		misc.SaveDroppedMessage(ctx, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, "TRANSACTION_PROCESSING_FAIL", msg.Value)
		return nil
	}

	return nil
}
