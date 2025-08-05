package database

import (
	"context"
	"prompt2video/internal/models"
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

