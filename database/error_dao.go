package database

import (
	"banking_ledger/logger"
	"banking_ledger/models"
	"banking_ledger/utils"
	"context"

	"fmt"
)

type errorDb struct{}

type errorDbInterface interface {
	SaveDroppedMessage(ctx context.Context, droppedAckEvent models.DroppedMessage) *models.ApplicationError
	ProcessErrorMessages(ctx context.Context, priority int, errorMessage string, moreInfo string) *models.ApplicationError
}

var ErDb errorDbInterface

func init() {
	ErDb = &errorDb{}
}

func (d *errorDb) SaveDroppedMessage(ctx context.Context, droppedMessage models.DroppedMessage) *models.ApplicationError {

	sqlStatement := `insert into kafka_topic_dropped_messages ("topic_name","error_type","kafka_message") values ($1,$2,$3)`
	_, err := dbPool.Exec(context.Background(), sqlStatement, droppedMessage.TopicName, droppedMessage.ErrorType, droppedMessage.KafkaMessage)

	if err != nil {
		errMsg := fmt.Sprintf("SaveDroppedMessage:Could not write to topic dropped message database table.Error:%s", err.Error())
		displayMsg := "Could not write to topic dropped message database table!"
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ctx, 3001, errMsg, displayMsg, nil)
		return appError
	}

	return nil

}

func (d *errorDb) ProcessErrorMessages(ctx context.Context, priority int, errorMessage string, additionalInfo string) *models.ApplicationError {

	sqlStatement := `insert into service_errors ("priority","error_message","additional_info") values ($1,$2,$3)`
	_, err := dbPool.Exec(context.Background(), sqlStatement, priority, errorMessage, additionalInfo)

	if err != nil {
		errMsg := fmt.Sprintf("ProcessErrorMessages:Could not write to errors table.Error:%s", err.Error())
		displayMsg := "Could not write to errors table!"
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ctx, 3002, errMsg, displayMsg, nil)
		return appError
	}

	return nil

}
