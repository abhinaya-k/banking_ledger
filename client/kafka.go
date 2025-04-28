package clients

import (
	"banking_ledger/config"
	"banking_ledger/handlers"
	"banking_ledger/logger"
	"banking_ledger/models"
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

type callbackFunctionWithMsg func(*kafka.Message) error

type ToKafkaMessage struct {
	Topic string
	Key   string
	Value []byte
}

var (
	ToKafkaChToTransactionProcessor = make(chan ToKafkaMessage)
)

func KafkaConsumer(kafkaConsumerGroup string, kafkaTopicName string, callbackFunction callbackFunctionWithMsg) {

	requestId := uuid.New().String()
	ctx := context.WithValue(context.Background(), models.CONTEXT_REQUEST_ID_KEY, requestId)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  config.AppConfig.KafkaBrokers,
		"group.id":           kafkaConsumerGroup,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
		"sasl.mechanisms":    "SCRAM-SHA-512",
		"security.protocol":  "sasl_ssl",
		"sasl.username":      config.AppConfig.KafkaUserName,
		"sasl.password":      config.AppConfig.KafkaPassword,
	})

	if err != nil {
		errorMsg := fmt.Sprintf("Kafka consumer connection error.Topic:%s,Error:%s!\n", kafkaTopicName, err.Error())
		logger.Log.Error(errorMsg)
		handlers.ProcessError(ctx, models.KAFKA_CONSUMER_ERROR, errorMsg, nil)
		panic(err)
	}

	err = c.SubscribeTopics([]string{kafkaTopicName}, nil)

	if err != nil {
		errorMsg := fmt.Sprintf("Kafka consumer subscribe error.Topic:%s,Error:%s!\n", kafkaTopicName, err.Error())
		logger.Log.Error(errorMsg)
		handlers.ProcessError(ctx, models.KAFKA_CONSUMER_ERROR, errorMsg, nil)
		panic(err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			err = callbackFunction(msg)
			if err == nil {
				c.CommitMessage(msg)
			} else {
				errorMsg := fmt.Sprintf("Kafka consumer commit error.Error:%s,Message:%s!\n", err.Error(), string(msg.Value))
				logger.Log.Error(errorMsg)
				handlers.ProcessError(ctx, models.KAFKA_CONSUMER_ERROR, errorMsg, msg.Value)
				c.Close()
				return
			}
		} else {
			//Here I am making the service panic and restart whenever consumer read error occurs
			errorMsg := fmt.Sprintf("Kafka consumer read error.Error:%s!\n", err.Error())
			logger.Log.Error(errorMsg)
			handlers.ProcessError(ctx, models.KAFKA_CONSUMER_ERROR, errorMsg, nil)
			panic(err)
		}
	}

}

func KafkaProducer(producerChannel <-chan ToKafkaMessage) {

	requestId := uuid.New().String()
	ctx := context.WithValue(context.Background(), models.CONTEXT_REQUEST_ID_KEY, requestId)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.AppConfig.KafkaBrokers,
		"sasl.mechanisms":   "SCRAM-SHA-512",
		"security.protocol": "sasl_ssl",
		"sasl.username":     config.AppConfig.KafkaUserName,
		"sasl.password":     config.AppConfig.KafkaPassword,
	})

	if err != nil {
		errorMsg := fmt.Sprintf("Kafka producer connection error.Error:%s!\n", err.Error())
		logger.Log.Error(errorMsg)
		handlers.ProcessError(ctx, models.KAFKA_PRODUCER_ERROR, errorMsg, nil)
		panic(err)
	}

	deliveryChan := make(chan kafka.Event)

	for {

		message := <-producerChannel

		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &(message.Topic), Partition: kafka.PartitionAny},
			Key:            []byte(message.Key),
			Value:          []byte(message.Value),
		}, deliveryChan)

		kafkaEvent := <-deliveryChan
		kafkaMessage := kafkaEvent.(*kafka.Message)

		if kafkaMessage.TopicPartition.Partition == -1 {
			errorMsg := fmt.Sprintf("Kafka producer or topic partition error.Topic:%s,Message:%s!\n", *kafkaMessage.TopicPartition.Topic, string(kafkaMessage.Value))
			logger.Log.Error(errorMsg)
			handlers.ProcessError(ctx, models.KAFKA_PRODUCER_ERROR, errorMsg, kafkaMessage.Value)
		}

		if err != nil {
			errorMsg := fmt.Sprintf("Kafka producer produce error.Error:%s!\n", err.Error())
			logger.Log.Error(errorMsg)
			handlers.ProcessError(ctx, models.KAFKA_PRODUCER_ERROR, errorMsg, nil)
		}

	}

}
