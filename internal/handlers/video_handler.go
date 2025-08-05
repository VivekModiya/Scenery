package handlers

import (
	"net/http"
	"strconv"
	"prompt2video/internal/models"
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
	var req models.CreateVideoJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Create video job
	job := &models.VideoJob{
		UserPrompt: req.Prompt,
		Subject:    req.Subject,
		Status:     models.JobStatusQueued,
	}

	if err := h.db.Create(job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create job",
		})
		return
	}

	// Queue the job for processing
	if err := h.jobService.QueueJob(job.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to queue job",
		})
		return
	}

	// Start processing asynchronously
	go h.processVideoJob(job.ID)

	c.JSON(http.StatusCreated, models.VideoJobResponse{
		ID:        job.ID,
		Status:    job.Status,
		CreatedAt: job.CreatedAt,
	})
}

// GetJobStatus returns the status of a video generation job
func (h *VideoHandler) GetJobStatus(c *gin.Context) {
	jobIDStr := c.Param("jobId")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid job ID",
		})
		return
	}

	var job models.VideoJob
	if err := h.db.First(&job, "id = ?", jobID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch job",
		})
		return
	}

	// Calculate progress based on status
	progress := 0
	switch job.Status {
	case models.JobStatusQueued:
		progress = 10
	case models.JobStatusProcessing:
		progress = 50
	case models.JobStatusCompleted:
		progress = 100
	case models.JobStatusFailed:
		progress = 0
	}

	response := models.JobStatusResponse{
		ID:          job.ID,
		Status:      job.Status,
		VideoURL:    job.VideoURL,
		Error:       job.ErrorMsg,
		Progress:    progress,
		CreatedAt:   job.CreatedAt,
		CompletedAt: job.CompletedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetVideo returns the video details
func (h *VideoHandler) GetVideo(c *gin.Context) {
	jobIDStr := c.Param("jobId")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid job ID",
		})
		return
	}

	var video models.GeneratedVideo
	if err := h.db.First(&video, "job_id = ?", jobID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Video not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch video",
		})
		return
	}

	c.JSON(http.StatusOK, video)
}

// ListVideos returns a paginated list of generated videos
func (h *VideoHandler) ListVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var videos []models.GeneratedVideo
	var total int64

	h.db.Model(&models.GeneratedVideo{}).Count(&total)
	
	if err := h.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch videos",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetJobDetails returns detailed information about a job
func (h *VideoHandler) GetJobDetails(c *gin.Context) {
	jobIDStr := c.Param("jobId")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid job ID",
		})
		return
	}

	var job models.VideoJob
	if err := h.db.First(&job, "id = ?", jobID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch job",
		})
		return
	}

	c.JSON(http.StatusOK, job)
}

// CancelJob cancels a video generation job
func (h *VideoHandler) CancelJob(c *gin.Context) {
	jobIDStr := c.Param("jobId")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid job ID",
		})
		return
	}

	var job models.VideoJob
	if err := h.db.First(&job, "id = ?", jobID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch job",
		})
		return
	}

	// Only allow cancellation if job is queued or processing
	if job.Status == models.JobStatusCompleted || job.Status == models.JobStatusFailed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot cancel completed or failed job",
		})
		return
	}

	// Update job status to failed
	job.Status = models.JobStatusFailed
	job.ErrorMsg = "Job cancelled by user"

	if err := h.db.Save(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to cancel job",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job cancelled successfully",
	})
}

// processVideoJob handles the async processing of a video job
func (h *VideoHandler) processVideoJob(jobID uuid.UUID) {
	// Update job status to processing
	h.db.Model(&models.VideoJob{}).Where("id = ?", jobID).Update("status", models.JobStatusProcessing)

	var job models.VideoJob
	if err := h.db.First(&job, "id = ?", jobID).Error; err != nil {
		return
	}

	// Step 1: Generate explanation and Manim code using LLM
	llmResponse, err := h.llmService.GenerateExplanation(job.UserPrompt, job.Subject)
	if err != nil {
		h.handleJobError(jobID, "LLM generation failed: "+err.Error())
		return
	}

	// Update job with LLM response
	job.LLMResponse = llmResponse.Explanation
	job.ManimCode = llmResponse.ManimCode
	h.db.Save(&job)

	// Step 2: Render video using Manim
	videoResponse, err := h.videoService.RenderVideo(models.VideoRenderRequest{
		JobID:     jobID,
		ManimCode: llmResponse.ManimCode,
		Prompt:    job.UserPrompt,
	})
	if err != nil {
		h.handleJobError(jobID, "Video rendering failed: "+err.Error())
		return
	}

	// Step 3: Update job as completed
	job.Status = models.JobStatusCompleted
	job.VideoURL = videoResponse.VideoURL
	now := job.UpdatedAt
	job.CompletedAt = &now
	h.db.Save(&job)

	// Step 4: Create GeneratedVideo record
	generatedVideo := &models.GeneratedVideo{
		JobID:      jobID,
		UserPrompt: job.UserPrompt,
		Subject:    job.Subject,
		VideoURL:   videoResponse.VideoURL,
		Duration:   15, // Default 15 seconds
		FileSize:   0,  // TODO: Calculate actual file size
	}
	h.db.Create(generatedVideo)
}

func (h *VideoHandler) handleJobError(jobID uuid.UUID, errorMsg string) {
	h.db.Model(&models.VideoJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
		"status":    models.JobStatusFailed,
		"error_msg": errorMsg,
	})
}
