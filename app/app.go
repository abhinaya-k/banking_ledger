package app

import (
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

	fmt.Println("starting")

	go func() {
		Router.Use(middleware.CorsMiddleware())
		if err := Router.Run(":8001"); err != nil {
			panic(err)
		}

		fmt.Println("app started on port:8001")
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

}
