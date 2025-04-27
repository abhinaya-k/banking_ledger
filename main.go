package main

import (
	"banking_ledger/app"
	"banking_ledger/config"
)

func main() {

	config.LoadEnv()
	config.ShowServiceInfo()
	app.StartApp()

}
