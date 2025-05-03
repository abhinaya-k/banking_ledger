package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	KafkaBrokers  string
	KafkaUserName string
	KafkaPassword string
	MongoHost     string
	MongoPort     int
	MongoDbName   string
}

var AppConfig Config

var (
	SERVICE_BASE_PATH                  string
	JWT_SECRET                         string
	TRANSACTION_PROCESSING_KAFKA_TOPIC string
	TRANSACTION_PROCESSING_KAFKA_CG    string
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Printf(".env file not found, reading configuration from environment! Error: %s", err.Error())
			panic(err)
		}
	}

}

func LoadEnv() {
	AppConfig = Config{
		ServerPort:    os.Getenv("SERVER_PORT"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		KafkaBrokers:  os.Getenv("KAFKA_BROKER"),
		KafkaUserName: os.Getenv("KAFKA_USERNAME"),
		KafkaPassword: os.Getenv("KAFKA_PASSWORD"),
		MongoHost:     os.Getenv("MONGO_HOST"),
		MongoPort:     getEnvAsInt("MONGO_PORT", 27017),
		MongoDbName:   os.Getenv("MONGO_DB_NAME"),
	}

	SERVICE_BASE_PATH = os.Getenv("SERVICE_BASE_PATH")
	JWT_SECRET = os.Getenv("JWT_SECRET")
	TRANSACTION_PROCESSING_KAFKA_TOPIC = os.Getenv("TRANSACTION_PROCESSING_KAFKA_TOPIC")
	TRANSACTION_PROCESSING_KAFKA_CG = os.Getenv("TRANSACTION_PROCESSING_KAFKA_CG")
}

// Helper function to read environment variable or fallback default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// If you want to load INT envs easily
func getEnvAsInt(name string, defaultVal int) int {
	valStr := getEnv(name, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func ShowServiceInfo() {

	serviceName := os.Getenv("SERVICE_NAME")
	fmt.Printf("Starting service %s!\n", serviceName)

	envName := os.Getenv("ENV")
	fmt.Printf("Environment is %s!\n", envName)

}
