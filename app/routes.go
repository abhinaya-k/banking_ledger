package app

import (
	"banking_ledger/config"
	"banking_ledger/handlers"

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

func SetupHealthRoute() {

	Router.GET(SERVICE_BASE_PATH+"/v1/health", handlers.GetHealth)
	Router.NoRoute(handlers.NoRoute)
}

func SetupUserRoute() {
	userRoutes.POST("/user/v1/register", handlers.RegisterUser)
}

func SetupCognitoProtectedRoutes() {

}
