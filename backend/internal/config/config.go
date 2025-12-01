package config

import (
	"fmt"
	"os"
)

type Config struct {
	MongoURI     string
	DatabaseName string
	JWTSecret    string
	
	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string
}

func Load() (*Config, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is required")
	}

	return &Config{
		MongoURI:     mongoURI,
		DatabaseName: getEnv("DATABASE_NAME", "telegraph"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPEmail:    getEnv("SMTP_EMAIL", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
