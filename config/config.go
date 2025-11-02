package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort           string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	RedisHost         string
	RedisPort         string
	RabbitMQURL       string
	ProductServiceURL string
	MockProduct       bool
}

func InitConfig() *Config {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	mockProductStr := os.Getenv("MOCK_PRODUCT_SERVICE")
	mockProduct, _ := strconv.ParseBool(mockProductStr)

	config := &Config{
		AppPort:           os.Getenv("APP_PORT"),
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		RedisHost:         os.Getenv("REDIS_HOST"),
		RedisPort:         os.Getenv("REDIS_PORT"),
		RabbitMQURL:       os.Getenv("RABBITMQ_URL"),
		ProductServiceURL: os.Getenv("PRODUCT_SERVICE_URL"),
		MockProduct:       mockProduct,
	}

	return config
}
