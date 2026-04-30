package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	BaseURL            string
	Port               string
	Private            bool
	APIKey             string
	MaxURLLength       int
	BlacklistedDomains []string
}

func Load() Config {
	loadEnv()

	return Config{
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		BaseURL:            strings.TrimRight(getEnv("BASE_URL", "http://localhost:4000"), "/"),
		Port:               getEnv("PORT", "4000"),
		Private:            getBoolEnv("PRIVATE", false),
		APIKey:             getEnv("API_KEY", ""),
		MaxURLLength:       getIntEnv("MAX_URL_LENGTH", 2048),
		BlacklistedDomains: getCSVEnv("BLACKLISTED_DOMAINS", "localhost,local,internal,metadata.google.internal"),
	}
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func getBoolEnv(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getIntEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}

func getCSVEnv(key string, fallback string) []string {
	value := getEnv(key, fallback)
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))

	for _, part := range parts {
		item := strings.ToLower(strings.TrimSpace(part))
		if item != "" {
			items = append(items, item)
		}
	}

	return items
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
