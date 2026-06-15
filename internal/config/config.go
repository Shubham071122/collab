package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	JWTSecret             string
	DatabaseURL           string
	ResendAPIKey          string
	ResendFromEmail       string
	AllowedOrigin         string
	RazorpayKeyID         string
	RazorpayKeySecret     string
	RazorpayWebhookSecret string
	RazorpaySilverPlanID  string
	RazorpayGoldPlanID    string
}

func LoadConfig() *Config {

	_ = godotenv.Load()

	return &Config{
		Port:                  os.Getenv("PORT"),
		JWTSecret:             os.Getenv("JWT_SECRET"),
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		ResendAPIKey:          os.Getenv("RESEND_API_KEY"),
		ResendFromEmail:       os.Getenv("RESEND_FROM_EMAIL"),
		AllowedOrigin:         os.Getenv("ALLOWED_ORIGIN"),
		RazorpayKeyID:         os.Getenv("RAZORPAY_KEY_ID"),
		RazorpayKeySecret:     os.Getenv("RAZORPAY_KEY_SECRET"),
		RazorpayWebhookSecret: os.Getenv("RAZORPAY_WEBHOOK_SECRET"),
		RazorpaySilverPlanID:  os.Getenv("RAZORPAY_SILVER_PLAN_ID"),
		RazorpayGoldPlanID:    os.Getenv("RAZORPAY_GOLD_PLAN_ID"),
	}
}