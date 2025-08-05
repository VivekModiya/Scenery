package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	RedisURL      string
	OpenAIAPIKey  string
	GeminiAPIKey  string
	VideoStorage  string
	MaxFileSize   int64
	WorkerCount   int
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://localhost/prompt2video?sslmode=disable"),
		VideoStorage:  getEnv("VIDEO_STORAGE", "./storage/videos"),
		MaxFileSize:   getEnvAsInt64("MAX_FILE_SIZE", 100*1024*1024), // 100MB
		WorkerCount:   getEnvAsInt("WORKER_COUNT", 3),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
