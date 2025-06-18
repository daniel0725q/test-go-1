package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/truora/microservice/internal/database"
	truoraHttp "github.com/truora/microservice/internal/delivery/http"
	"github.com/truora/microservice/internal/repository"
	"github.com/truora/microservice/internal/usecase"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	ExternalAPI struct {
		BaseURL string `yaml:"base_url"`
		Timeout int    `yaml:"timeout"`
		Token   string `yaml:"token"`
	} `yaml:"external_api"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Load configuration
	configFile, err := os.ReadFile("config/config.yml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(configFile, config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Initialize database connection
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
		config.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(sqlDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	stockRatingRepo := repository.NewStockRatingRepository(db)
	jobRepo := repository.NewJobRepository(db)
	externalAPIRepo := repository.NewExternalAPIRepository(
		config.ExternalAPI.BaseURL,
		time.Duration(config.ExternalAPI.Timeout)*time.Second,
		config.ExternalAPI.Token,
	)

	// Initialize service
	stockRatingSvc := usecase.NewStockRatingService(stockRatingRepo, jobRepo, externalAPIRepo)

	// Initialize handler
	handler := truoraHttp.NewHandler(stockRatingSvc)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// Register routes
	handler.RegisterRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
