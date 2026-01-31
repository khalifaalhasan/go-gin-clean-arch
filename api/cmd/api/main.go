package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/config"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/database"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/handler"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/middleware"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/service"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// 2. Konek Database
	conn, err := database.NewConnection(cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	
	defer conn.Close()


	store := database.New(conn)

	server := NewServer(store)

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	

	log.Println("Database connected successfully!")
	
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	turnstileSecret := os.Getenv("TURNSTILE_SECRET_KEY")
	if turnstileSecret == "" {
		turnstileSecret = "1x0000000000000000000000000000000AA" // 
	}


	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	
	if _, err := rdb.Ping(rdb.Context()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	
	secService := service.NewSecurityService(turnstileSecret)
	commentHandler := handler.NewCommentHandler(secService)


	r := gin.Default() 
	


r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"}, // Izinkan Next.js
    AllowMethods:     []string{"POST", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "X-RateLimit-Remaining"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
}))
	
	


	api := r.Group("/api/v1")
	{
		
		api.POST("/comments", middleware.RateLimiter(rdb), commentHandler.PostComment)
	}

	
	log.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}