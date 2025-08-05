package handlers

import (
	"prompt2video/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JobServiceInterface defines the contract for job service
type JobServiceInterface interface {
	QueueJob(jobID uuid.UUID) error
}

type VideoHandler struct {
	db           *gorm.DB
	llmService   *services.LLMService
	videoService *services.VideoService
	jobService   JobServiceInterface
}

func NewVideoHandler(
	db *gorm.DB,
	llmService *services.LLMService,
	videoService *services.VideoService,
	jobService JobServiceInterface,
) *VideoHandler {
	return &VideoHandler{
		db:           db,
		llmService:   llmService,
		videoService: videoService,
		jobService:   jobService,
	}
}

// GenerateVideo creates a new video generation job
func (h *VideoHandler) GenerateVideo(c *gin.Context) {
// step 1 get script/explaniation abt video 

//	llmExplaination, err := h.llmService.GenerateExplanation(userPrompt)

// step 2 llmExplainaton to get mainm code 

//	llmMainm, err := h.llmService.GenerateMainm(userPrompt)

// step 3 use lllmanim code to render video
//	videoResponse, err := h.videoService.RenderVideo(llmMainm)
}








