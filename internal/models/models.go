package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VideoJob represents a video generation job
type VideoJob struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserPrompt  string    `gorm:"not null" json:"user_prompt"`
	Subject     string    `gorm:"default:general" json:"subject"`
	Status      JobStatus `gorm:"default:queued" json:"status"`
	LLMResponse string    `json:"llm_response,omitempty"`
	ManimCode   string    `json:"manim_code,omitempty"`
	VideoURL    string    `json:"video_url,omitempty"`
	ErrorMsg    string    `json:"error_message,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// JobStatus represents the status of a video generation job
type JobStatus string

const (
	JobStatusQueued     JobStatus = "queued"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

// GeneratedVideo represents a successfully generated video
type GeneratedVideo struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	JobID       uuid.UUID `gorm:"type:uuid;not null" json:"job_id"`
	UserPrompt  string    `gorm:"not null" json:"user_prompt"`
	Subject     string    `json:"subject"`
	VideoURL    string    `gorm:"not null" json:"video_url"`
	Duration    int       `json:"duration"` // in seconds
	FileSize    int64     `json:"file_size"` // in bytes
	CreatedAt   time.Time `json:"created_at"`
}

// CreateVideoJobRequest represents the request to create a video job
type CreateVideoJobRequest struct {
	Prompt  string `json:"prompt" binding:"required,min=3,max=500"`
	Subject string `json:"subject,omitempty"`
}

// VideoJobResponse represents the response when creating/getting a video job
type VideoJobResponse struct {
	ID        uuid.UUID `json:"id"`
	Status    JobStatus `json:"status"`
	VideoURL  string    `json:"video_url,omitempty"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// JobStatusResponse represents the response for job status checks
type JobStatusResponse struct {
	ID          uuid.UUID  `json:"id"`
	Status      JobStatus  `json:"status"`
	VideoURL    string     `json:"video_url,omitempty"`
	Error       string     `json:"error,omitempty"`
	Progress    int        `json:"progress"` // 0-100
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// LLMRequest represents a request to the LLM service
type LLMRequest struct {
	Prompt  string `json:"prompt"`
	Subject string `json:"subject"`
}

// LLMResponse represents the response from the LLM service
type LLMResponse struct {
	Explanation string `json:"explanation"`
	ManimCode   string `json:"manim_code"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
}

// VideoRenderRequest represents a request to render a video
type VideoRenderRequest struct {
	JobID     uuid.UUID `json:"job_id"`
	ManimCode string    `json:"manim_code"`
	Prompt    string    `json:"prompt"`
}

// VideoRenderResponse represents the response from video rendering
type VideoRenderResponse struct {
	Success  bool   `json:"success"`
	VideoURL string `json:"video_url,omitempty"`
	Error    string `json:"error,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (vj *VideoJob) BeforeCreate(tx *gorm.DB) error {
	if vj.ID == uuid.Nil {
		vj.ID = uuid.New()
	}
	return nil
}

func (gv *GeneratedVideo) BeforeCreate(tx *gorm.DB) error {
	if gv.ID == uuid.Nil {
		gv.ID = uuid.New()
	}
	return nil
}
