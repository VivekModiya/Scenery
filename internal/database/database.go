package database

import (
	"context"
	"prompt2video/internal/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate models
	err = db.AutoMigrate(
		&models.VideoJob{},
		&models.GeneratedVideo{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitializeRedis(redisURL string) *redis.Client {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		// Fallback to default configuration
		opts = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	client := redis.NewClient(opts)

	// Test connection
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return client
}
