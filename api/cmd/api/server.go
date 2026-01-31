package main

import (
	"github.com/gin-gonic/gin"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/database"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/handler"
)

type Server struct {
	store *database.Queries
	router *gin.Engine
}


func NewServer(db *database.Queries) *Server {
	server := &Server{store: db}
	router := gin.Default()

	userHandler := handler.NewUserHandler(server.store)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", userHandler.CreateUser)
	}

	server.router = router
	return server
}

func (s *Server) Start(address string) error{
	return s.router.Run(address)
}