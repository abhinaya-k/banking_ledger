package database

import (
	"banking_ledger/config"
	"banking_ledger/logger"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func InitializeDatabasePool() error {

	var err error

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.AppConfig.DBUser, config.AppConfig.DBPassword, config.AppConfig.DBHost, config.AppConfig.DBPort, config.AppConfig.DBName)
	fmt.Println(databaseUrl)

	dbconfig, _ := pgxpool.ParseConfig(databaseUrl)
	dbPool, err = pgxpool.NewWithConfig(context.Background(), dbconfig)
	if err != nil {
		errorMsg := fmt.Sprintf("Cannot connect to database %s.Error:%s!\n", databaseUrl, err.Error())
		logger.Log.Error(errorMsg)
		return err
	}

	err = PingDatabasePool()
	if err != nil {
		errMsg := fmt.Sprintf("Unable to ping database: %v\n", err)
		logger.Log.Error(errMsg)
		return err
	}

	logger.Log.Info("Successfully connected to database")

	return nil

}

func PingDatabasePool() error {

	return dbPool.Ping(context.Background())

}

func CloseDatabasePool() {

	dbPool.Close()

}
