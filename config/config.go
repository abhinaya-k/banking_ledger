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
	KafkaBrokers  []string
	KafkaTopic    string
	KafkaUserName string
	KafkaPassword string
}

var AppConfig Config

var (
	SERVICE_BASE_PATH string
	JWT_SECRET        string
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
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "password"),
		DBName:        getEnv("DB_NAME", "ledgerdb"),
		KafkaBrokers:  []string{getEnv("KAFKA_BROKER", "localhost:9092")},
		KafkaTopic:    getEnv("KAFKA_TOPIC", "transaction_events"),
		KafkaUserName: getEnv("KAFKA_USERNAME", "admin"),
		KafkaPassword: getEnv("KAFKA_PASSWORD", "admin"),
	}

	SERVICE_BASE_PATH = getEnv("SERVICE_BASE_PATH", "/bankingLedger")
	JWT_SECRET = os.Getenv("JWT_SECRET")
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
