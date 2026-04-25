package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL      string
	RedisAddr        string
	RedisPass        string
	Port             string
	TelegramBotToken string
}

func Load() *Config {
	godotenv.Load("../.env")
	return &Config{
		PostgresURL:      os.Getenv("POSTGRES_URL"),
		RedisAddr:        "localhost:6379",
		RedisPass:        os.Getenv("REDIS_PASSWORD"),
		Port:             getEnvOrDefault("BACKEND_PORT", "8080"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
