package config

import (
	"os"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	// Save original env values
	originalPort := os.Getenv("PORT")
	originalDBURL := os.Getenv("DATABASE_URL")
	
	// Set test values
	os.Setenv("PORT", "9000")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	
	// Load config
	cfg := Load()
	
	// Test values
	if cfg.Port != "9000" {
		t.Errorf("Expected Port to be '9000', got %s", cfg.Port)
	}
	
	if cfg.DatabaseURL != "postgres://test:test@localhost:5432/testdb" {
		t.Errorf("Expected DatabaseURL to be test value, got %s", cfg.DatabaseURL)
	}
	
	// Restore original values
	if originalPort != "" {
		os.Setenv("PORT", originalPort)
	} else {
		os.Unsetenv("PORT")
	}
	
	if originalDBURL != "" {
		os.Setenv("DATABASE_URL", originalDBURL)
	} else {
		os.Unsetenv("DATABASE_URL")
	}
}

func TestConfigDefaults(t *testing.T) {
	// Clear env vars to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("WORKER_COUNT")
	
	cfg := Load()
	
	if cfg.Port != "8080" {
		t.Errorf("Expected default Port to be '8080', got %s", cfg.Port)
	}
	
	if cfg.WorkerCount != 3 {
		t.Errorf("Expected default WorkerCount to be 3, got %d", cfg.WorkerCount)
	}
	
	if cfg.MaxFileSize != 100*1024*1024 {
		t.Errorf("Expected default MaxFileSize to be 100MB, got %d", cfg.MaxFileSize)
	}
}
