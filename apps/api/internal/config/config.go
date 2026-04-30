package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	BaseURL     string
	Port        string
}

func Load() Config {
	loadEnv()

	return Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		BaseURL:     strings.TrimRight(getEnv("BASE_URL", "http://localhost:4000"), "/"),
		Port:        getEnv("PORT", "4000"),
	}
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func loadEnv() {
	candidates := []string{
		".env",
		"../../.env",
		"apps/api/.env",
	}

	for _, filename := range candidates {
		if _, err := os.Stat(filename); err != nil {
			if !os.IsNotExist(err) {
				log.Printf("could not inspect %s: %v", filename, err)
			}
			continue
		}

		if err := godotenv.Load(filename); err != nil {
			log.Printf("could not load %s: %v", filename, err)
		}
	}
}
