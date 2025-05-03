package config

import (
	env "admin-phone-shop-api/pkg/utils/env"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppHost string
	AppPort int
}

func NewConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	host := os.Getenv("API_HOST")
	port := env.GetenvInt("API_PORT", 8889)

	return &AppConfig{
		AppHost: host,
		AppPort: port,
	}
}
