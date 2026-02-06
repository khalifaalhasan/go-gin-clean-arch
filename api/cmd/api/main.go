package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/config"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/database"
	
	"github.com/khalifaalhasan/go-gin-clean-arch/api"
)
	

func main() {
	// 1. Load config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// 2. Konek Database
	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	// Tes Koneksi DB
	if err:= conn.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}
	log.Println("Database Connected successfully!")

	// Initialize database store (SQLC)
	store := database.New(conn)

	// Konek Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB : 0,
	})

	// tes koneksi redis
	if _, err := rdb.Ping(rdb.Context()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Redis connected successfully!")

	// Buat Server 
	server, err := api.NewServer(cfg, store, conn, rdb)
	if err != nil{
		log.Fatal("cannot create server:", err)
	}

	// Run Server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	
}	

