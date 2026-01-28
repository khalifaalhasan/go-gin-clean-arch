package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/handler"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/middleware"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/service"
	"github.com/gin-contrib/cors"
)

func main() {
	// 1. Setup Environment (Mocking ENV for demonstration)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	turnstileSecret := os.Getenv("TURNSTILE_SECRET_KEY")
	if turnstileSecret == "" {
		turnstileSecret = "1x0000000000000000000000000000000AA" // Mock Secret Key Cloudflare
	}

	// 2. Setup Infrastructure (Redis)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	// Ping Redis to ensure connection
	if _, err := rdb.Ping(rdb.Context()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	// 3. Dependency Injection
	secService := service.NewSecurityService(turnstileSecret)
	commentHandler := handler.NewCommentHandler(secService)

	// 4. Gin Setup
	r := gin.Default() // Default includes Logger and Recovery middleware

	// ... import "github.com/gin-contrib/cors" (Jangan lupa go get dulu)

// Di dalam main(), setelah r := gin.Default()
r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"}, // Izinkan Next.js
    AllowMethods:     []string{"POST", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "X-RateLimit-Remaining"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
}))
	
	// Trusted Proxies configuration (PENTING jika di belakang Nginx/Cloudflare)
	// r.SetTrustedProxies([]string{"127.0.0.1"}) 

	// 5. Routes Definition
	api := r.Group("/api/v1")
	{
		// Terapkan Rate Limiter hanya pada endpoint yang rawan (Input Komentar)
		api.POST("/comments", middleware.RateLimiter(rdb), commentHandler.PostComment)
	}

	// 6. Run Server
	log.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}