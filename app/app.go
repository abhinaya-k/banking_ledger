package app

import (
	"banking_ledger/config"
	"banking_ledger/database"
	"banking_ledger/middleware"
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
	}()

	SetupHealthRoute()

	if err := database.InitializeDatabasePool(); err != nil {
		panic(err)
	}

	SERVER_PORT := fmt.Sprintf(":%s", config.AppConfig.ServerPort)

	go func() {
		Router.Use(middleware.CorsMiddleware())
		if err := Router.Run(SERVER_PORT); err != nil {
			panic(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

}
