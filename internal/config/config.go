package config

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	JWTSecret       string
	DatabaseURL     string
	ResendAPIKey    string
	ResendFromEmail string
}

func LoadConfig() *Config {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		Port:            os.Getenv("PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		ResendAPIKey:    os.Getenv("RESEND_API_KEY"),
		ResendFromEmail: os.Getenv("RESEND_FROM_EMAIL"),
	}
}