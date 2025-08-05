package api

import (
	"net/http"
	"prompt2video/internal/config"
	"prompt2video/internal/handlers"
	"prompt2video/internal/middleware"
	"prompt2video/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Server struct {
	config       *config.Config
	db           *gorm.DB
	redis        *redis.Client
	llmService   *services.LLMService
	videoService *services.VideoService
	jobService   handlers.JobServiceInterface
	router       *gin.Engine
}

func NewServer(
	cfg *config.Config,
	db *gorm.DB,
	redisClient *redis.Client,
	llmService *services.LLMService,
	videoService *services.VideoService,
	jobService handlers.JobServiceInterface,
) *Server {
	server := &Server{
		config:       cfg,
		db:           db,
		redis:        redisClient,
		llmService:   llmService,
		videoService: videoService,
		jobService:   jobService,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	
	s.router = gin.New()

	// Middleware
	s.router.Use(middleware.CORS())

	// Health check

	// Initialize handlers
	videoHandler := handlers.NewVideoHandler(s.db, s.llmService, s.videoService, s.jobService)

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Video generation endpoints
		api.POST("/videos/generate", videoHandler.GenerateVideo)
		
	}

	// Serve static files (videos)
	s.router.Static("/storage", s.config.VideoStorage)
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.config.Port)
}
