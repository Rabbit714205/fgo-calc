package main

import (
	"fgo-calc-backend/internal/config"
	"fgo-calc-backend/internal/handler"
	"fgo-calc-backend/internal/repository"
	"fgo-calc-backend/internal/service"
	"flag"
	"log"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.vercel.json", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Vercel provides the PORT environment variable.
	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = ":" + port
	}

	repo, err := repository.NewRepository(cfg.DataDir)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	svc := service.NewCalculatorService(repo)

	h, err := handler.NewHandler(repo, svc, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize handler: %v", err)
	}

	r := gin.Default()

	r.Use(gzip.Gzip(
		gzip.DefaultCompression,
		gzip.WithExcludedExtensions([]string{
			".png", ".gif", ".jpeg", ".jpg", ".webp",
		}),
	))

	h.Register(r)

	log.Printf("Server starting on %s", cfg.Port)

	if err := r.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
