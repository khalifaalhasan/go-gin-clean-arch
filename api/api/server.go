package api

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/config"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/database"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/handler"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/service"
)

type Server struct {
	config config.Config 
	store *database.Queries
	db    *sql.DB
	rdb *redis.Client
	router *gin.Engine
}

// NewServer creates a new HTTP server
func NewServer(cfg config.Config, store *database.Queries, db *sql.DB, rdb *redis.Client) (*Server, error) {
	server := &Server{
		config: cfg, 
		store: store, 
		db: db, 
		rdb: rdb,
	}

	server.setupRouter()
	return server, nil
}

// setupRouter sets up the HTTP routes for the server
func (server *Server) setupRouter() {
	router := gin.Default()

	// setup CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins : []string{"http://localhost:3000"},
		AllowMethods : []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "Idempotency-Key"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// services
	secService := service.NewSecurityService(server.config.TurnstileSecret)

	// handler
	userHandler := handler.NewUserHandler(server.store)
	orderHandler := handler.NewOrderHandler(server.store, server.db)
	commentHandler := handler.NewCommentHandler(secService)

	// Routing Group
	api := router.Group("/api/v1")
	{
		//user routes
		api.POST("/users", userHandler.CreateUser)
		// order routes
		api.POST("/orders", orderHandler.CreateOrder)
		// ada middleware redis rate limiter
		api.POST("/comments", commentHandler.PostComment)
	}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
