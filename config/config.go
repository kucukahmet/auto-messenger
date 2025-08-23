package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTP_PORT         string
	DB_URI            string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	POSTGRES_PORT     string
	REDIS_PORT        string
	PERIOD            int
	BATCH_SIZE        int
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Load() (*Config, error) {
	godotenv.Load()

	httpPort := getEnv("HTTP_PORT", "8080")
	postgresUser := getEnv("POSTGRES_USER", "auto")
	postgresPassword := getEnv("POSTGRES_PASSWORD", "messager")
	postgresDb := getEnv("POSTGRES_DB", "automessager")
	postgresPort := getEnv("POSTGRES_PORT", "5432")
	postgresHost := getEnv("POSTGRES_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")

	period, _ := strconv.Atoi(getEnv("PERIOD", "120"))
	batchSize, _ := strconv.Atoi(getEnv("BATCH_SIZE", "2"))

	dbUri := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		postgresUser,
		postgresPassword,
		postgresHost,
		postgresPort,
		postgresDb,
	)

	return &Config{
		HTTP_PORT:         httpPort,
		DB_URI:            dbUri,
		POSTGRES_USER:     postgresUser,
		POSTGRES_PASSWORD: postgresPassword,
		POSTGRES_DB:       postgresDb,
		POSTGRES_PORT:     postgresPort,
		REDIS_PORT:        redisPort,
		PERIOD:            period,
		BATCH_SIZE:        batchSize,
	}, nil
}
