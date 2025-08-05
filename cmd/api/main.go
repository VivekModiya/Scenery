package main

import (
	"log"

	"prompt2video/internal/api"
	"prompt2video/internal/config"
	"prompt2video/internal/database"
	"prompt2video/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Redis
	redisClient := database.InitializeRedis(cfg.RedisURL)

	// Initialize services
	llmService := services.NewLLMService(cfg.OpenAIAPIKey, cfg.GeminiAPIKey)
	videoService := services.NewVideoService()
	jobService := services.NewJobService(redisClient)

	// Initialize API server
	server := api.NewServer(cfg, db, redisClient, llmService, videoService, jobService)

	// Start the server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
