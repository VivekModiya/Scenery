package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// RateLimit middleware using Redis
func RateLimit(redisClient *redis.Client) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Skip rate limiting if Redis client is nil (for testing)
		if redisClient == nil {
			c.Next()
			return
		}
		
		// Get client IP
		clientIP := c.ClientIP()
		
		// Create rate limit key
		key := fmt.Sprintf("rate_limit:%s", clientIP)
		
		ctx := context.Background()
		
		// Check current count
		count, err := redisClient.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			// If Redis is down, allow the request
			c.Next()
			return
		}
		
		// Rate limit: 100 requests per hour
		limit := 100
		window := time.Hour
		
		if count >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			c.Abort()
			return
		}
		
		// Increment counter
		pipe := redisClient.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)
		_, err = pipe.Exec(ctx)
		
		if err != nil {
			// If Redis is down, allow the request
			c.Next()
			return
		}
		
		c.Next()
	})
}

// RequestLogger middleware for detailed logging
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// Recovery middleware with custom error handling
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"message": err,
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
