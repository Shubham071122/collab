package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	JWTSecret       string
	DatabaseURL     string
	ResendAPIKey    string
	ResendFromEmail string
	AllowedOrigin   string
}

func LoadConfig() *Config {

	_ = godotenv.Load()

	return &Config{
		Port:            os.Getenv("PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		ResendAPIKey:    os.Getenv("RESEND_API_KEY"),
		ResendFromEmail: os.Getenv("RESEND_FROM_EMAIL"),
		AllowedOrigin:   os.Getenv("ALLOWED_ORIGIN"),
	}
}