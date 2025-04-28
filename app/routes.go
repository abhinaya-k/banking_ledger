package app

import (
	"banking_ledger/config"
	"banking_ledger/handlers"
	"banking_ledger/middleware"

	"github.com/gin-gonic/gin"
)

var (
	SERVICE_BASE_PATH      string
	cognitoProtectedRoutes *gin.RouterGroup
	userRoutes             *gin.RouterGroup
)

func init() {
	SERVICE_BASE_PATH = config.SERVICE_BASE_PATH
	cognitoProtectedRoutes = Router.Group(SERVICE_BASE_PATH)
	userRoutes = Router.Group(SERVICE_BASE_PATH)
}

func SetupRoutesMiddleware() {

	cognitoProtectedRoutes.Use(middleware.AuthTokenMiddleware())
}

func SetupHealthRoute() {

	Router.GET(SERVICE_BASE_PATH+"/v1/health", handlers.GetHealth)
	Router.NoRoute(handlers.NoRoute)
}

func SetupUserRoute() {
	userRoutes.POST("/user/v1/register", handlers.RegisterUser)
	userRoutes.POST("/user/v1/login", handlers.UserLogin)
}

func SetupCognitoProtectedRoutes() {

	cognitoProtectedRoutes.POST("/v1/account", handlers.CreateAccount)
	cognitoProtectedRoutes.PATCH("/account/v1/transaction", handlers.FundTransaction)
	cognitoProtectedRoutes.POST("/v1/transactionHistoty", handlers.GetTransactionHistory)

}
