// db/mongo.go
package database

import (
	"banking_ledger/config"
	"banking_ledger/logger"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient  *mongo.Client
	databaseName = config.AppConfig.MongoDbName
)

func InitMongoDB() error {
	uri := fmt.Sprintf("mongodb://%s:%d", config.AppConfig.MongoHost, config.AppConfig.MongoPort)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		errMsg := fmt.Sprintf("Error connecting to MongoDB: %v", err)
		logger.Log.Error(errMsg)
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		errMsg := fmt.Sprintf("MongoDB ping failed: %v", err)
		logger.Log.Error(errMsg)
		return err
	}

	MongoClient = client
	logger.Log.Info("Connected to MongoDB successfully")

	return nil
}

func DisconnectMongoDB() {
	if MongoClient == nil {
		return // No client, nothing to disconnect
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := MongoClient.Disconnect(ctx); err != nil {
		errMsg := fmt.Sprintf("Error disconnecting MongoDB: %v", err)
		logger.Log.Error(errMsg)
	} else {
		logger.Log.Info("MongoDB disconnected successfully")
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	if MongoClient == nil {
		panic("MongoDB not initialized")
	}

	return MongoClient.Database(databaseName).Collection(collectionName)
}
