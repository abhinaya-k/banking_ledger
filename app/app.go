package app

import (
	"banking_ledger/clients"
	"banking_ledger/config"
	"banking_ledger/database"
	"banking_ledger/middleware"
	"banking_ledger/services"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
)

func init() {

	gin.SetMode(gin.ReleaseMode)

	Router = gin.New()

}

func StartApp() {

	defer func() {
		database.CloseDatabasePool()
		database.DisconnectMongoDB()
	}()

	SetupHealthRoute()
	SetupRoutesMiddleware()
	SetupUserRoute()
	SetupCognitoProtectedRoutes()

	if err := database.InitializeDatabasePool(); err != nil {
		panic(err)
	}

	if err := database.InitMongoDB(); err != nil {
		panic(err)
	}

	SERVER_PORT := fmt.Sprintf(":%s", config.AppConfig.ServerPort)

	go func() {
		Router.Use(middleware.CorsMiddleware())
		if err := Router.Run(SERVER_PORT); err != nil {
			panic(err)
		}
	}()

	go clients.KafkaProducer(clients.ToKafkaChToTransactionProcessor)

	go clients.KafkaConsumer(config.TRANSACTION_PROCESSING_KAFKA_CG, config.TRANSACTION_PROCESSING_KAFKA_TOPIC, services.KafkaConsumerProcessTransactions)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

}
