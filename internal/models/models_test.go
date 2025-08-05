package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestVideoJobModel(t *testing.T) {
	job := VideoJob{
		UserPrompt: "Test prompt",
		Subject:    "mathematics",
		Status:     JobStatusQueued,
	}

	if job.UserPrompt != "Test prompt" {
		t.Errorf("Expected UserPrompt to be 'Test prompt', got %s", job.UserPrompt)
	}

	if job.Status != JobStatusQueued {
		t.Errorf("Expected Status to be JobStatusQueued, got %s", job.Status)
	}
}

func TestJobStatusValues(t *testing.T) {
	statuses := []JobStatus{
		JobStatusQueued,
		JobStatusProcessing,
		JobStatusCompleted,
		JobStatusFailed,
	}

	expectedValues := []string{
		"queued",
		"processing",
		"completed",
		"failed",
	}

	for i, status := range statuses {
		if string(status) != expectedValues[i] {
			t.Errorf("Expected status %d to be %s, got %s", i, expectedValues[i], string(status))
		}
	}
}

func TestCreateVideoJobRequest(t *testing.T) {
	req := CreateVideoJobRequest{
		Prompt:  "Explain quantum physics",
		Subject: "physics",
	}

	if req.Prompt != "Explain quantum physics" {
		t.Errorf("Expected Prompt to be 'Explain quantum physics', got %s", req.Prompt)
	}

	if req.Subject != "physics" {
		t.Errorf("Expected Subject to be 'physics', got %s", req.Subject)
	}
}

func TestGeneratedVideoModel(t *testing.T) {
	jobID := uuid.New()
	video := GeneratedVideo{
		JobID:      jobID,
		UserPrompt: "Test video",
		Subject:    "test",
		VideoURL:   "/storage/test.mp4",
		Duration:   15,
		FileSize:   1024000,
		CreatedAt:  time.Now(),
	}

	if video.JobID != jobID {
		t.Errorf("Expected JobID to match, got different values")
	}

	if video.Duration != 15 {
		t.Errorf("Expected Duration to be 15, got %d", video.Duration)
	}

	if video.FileSize != 1024000 {
		t.Errorf("Expected FileSize to be 1024000, got %d", video.FileSize)
	}
}
