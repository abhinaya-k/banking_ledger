package app

import (
	"banking_ledger/config"
	"banking_ledger/handlers"
)

func SetupHealthRoute() {
	Router.GET(config.SERVICE_BASE_PATH+"/v1/health", handlers.GetHealth)

	Router.NoRoute(handlers.NoRoute)
}
