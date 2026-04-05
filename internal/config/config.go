package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DBUrl string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port: getEnv("PORT", "8080"),
		DBUrl: getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}