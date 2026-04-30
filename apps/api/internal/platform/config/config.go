package config

import "os"

type AppConfig struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

var Current AppConfig

func Load() {
	Current = AppConfig{
		Port:        getEnv("PORT", "4000"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
