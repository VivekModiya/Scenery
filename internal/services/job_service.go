package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type JobService struct {
	redis *redis.Client
}

type JobData struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewJobService(redisClient *redis.Client) *JobService {
	return &JobService{
		redis: redisClient,
	}
}

// QueueJob adds a job to the processing queue
func (s *JobService) QueueJob(jobID uuid.UUID) error {
	ctx := context.Background()
	
	jobData := JobData{
		ID:        jobID,
		Status:    "queued",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Serialize job data
	data, err := json.Marshal(jobData)
	if err != nil {
		return fmt.Errorf("failed to marshal job data: %w", err)
	}

	// Add to queue
	queueKey := "video_jobs:queue"
	if err := s.redis.LPush(ctx, queueKey, string(data)).Err(); err != nil {
		return fmt.Errorf("failed to queue job: %w", err)
	}

	// Store job status
	statusKey := fmt.Sprintf("video_jobs:status:%s", jobID.String())
	if err := s.redis.Set(ctx, statusKey, string(data), 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to set job status: %w", err)
	}

	return nil
}

// GetJobStatus retrieves the current status of a job
func (s *JobService) GetJobStatus(jobID uuid.UUID) (*JobData, error) {
	ctx := context.Background()
	statusKey := fmt.Sprintf("video_jobs:status:%s", jobID.String())

	data, err := s.redis.Get(ctx, statusKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job status: %w", err)
	}

	var jobData JobData
	if err := json.Unmarshal([]byte(data), &jobData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job data: %w", err)
	}

	return &jobData, nil
}

// UpdateJobStatus updates the status of a job
func (s *JobService) UpdateJobStatus(jobID uuid.UUID, status string) error {
	ctx := context.Background()
	
	// Get current job data
	currentJob, err := s.GetJobStatus(jobID)
	if err != nil {
		// If job doesn't exist, create new one
		currentJob = &JobData{
			ID:        jobID,
			CreatedAt: time.Now(),
		}
	}

	// Update status and timestamp
	currentJob.Status = status
	currentJob.UpdatedAt = time.Now()

	// Serialize updated data
	data, err := json.Marshal(currentJob)
	if err != nil {
		return fmt.Errorf("failed to marshal job data: %w", err)
	}

	// Update status in Redis
	statusKey := fmt.Sprintf("video_jobs:status:%s", jobID.String())
	if err := s.redis.Set(ctx, statusKey, string(data), 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}

// DequeueJob removes and returns the next job from the queue
func (s *JobService) DequeueJob() (*JobData, error) {
	ctx := context.Background()
	queueKey := "video_jobs:queue"

	// Block for up to 30 seconds waiting for a job
	result, err := s.redis.BRPop(ctx, 30*time.Second, queueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No jobs available
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid queue result")
	}

	var jobData JobData
	if err := json.Unmarshal([]byte(result[1]), &jobData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job data: %w", err)
	}

	return &jobData, nil
}

// GetQueueLength returns the current number of jobs in the queue
func (s *JobService) GetQueueLength() (int64, error) {
	ctx := context.Background()
	queueKey := "video_jobs:queue"

	length, err := s.redis.LLen(ctx, queueKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}

	return length, nil
}

// RemoveJob removes a job from the system
func (s *JobService) RemoveJob(jobID uuid.UUID) error {
	ctx := context.Background()
	
	// Remove job status
	statusKey := fmt.Sprintf("video_jobs:status:%s", jobID.String())
	if err := s.redis.Del(ctx, statusKey).Err(); err != nil {
		return fmt.Errorf("failed to remove job status: %w", err)
	}

	// Note: We don't remove from queue as it's processed sequentially
	// The job will be skipped if it's already processed when dequeued

	return nil
}

// StartWorkers starts background workers to process video jobs
func (s *JobService) StartWorkers(workerCount int, processFunc func(uuid.UUID) error) {
	for i := 0; i < workerCount; i++ {
		go s.worker(i, processFunc)
	}
}

func (s *JobService) worker(workerID int, processFunc func(uuid.UUID) error) {
	fmt.Printf("Worker %d started\n", workerID)
	
	for {
		// Try to get a job from the queue
		job, err := s.DequeueJob()
		if err != nil {
			fmt.Printf("Worker %d: Error dequeuing job: %v\n", workerID, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if job == nil {
			// No jobs available, continue polling
			continue
		}

		fmt.Printf("Worker %d: Processing job %s\n", workerID, job.ID.String())

		// Update job status to processing
		if err := s.UpdateJobStatus(job.ID, "processing"); err != nil {
			fmt.Printf("Worker %d: Error updating job status: %v\n", workerID, err)
		}

		// Process the job
		if err := processFunc(job.ID); err != nil {
			fmt.Printf("Worker %d: Error processing job %s: %v\n", workerID, job.ID.String(), err)
			s.UpdateJobStatus(job.ID, "failed")
		} else {
			fmt.Printf("Worker %d: Successfully processed job %s\n", workerID, job.ID.String())
			s.UpdateJobStatus(job.ID, "completed")
		}
	}
}
