package main

import (
	"context"
	"mini-trh-backend/internal/logger"
	"mini-trh-backend/pkg/api/servers"
	"mini-trh-backend/pkg/infrastructure/postgres"
	"mini-trh-backend/pkg/infrastructure/postgres/repositories"
	"mini-trh-backend/pkg/infrastructure/postgres/schemas"
	"mini-trh-backend/pkg/services"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Log.Info("Starting Mini TRH Backend...")

	if err := godotenv.Load(); err != nil {
		logger.Log.Warn("No .env file found, using environment variables")
	}

	db, err := postgres.ConnectDatabase()
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := schemas.RunMigrations(db); err != nil {
		logger.Log.Fatal("Failed to run migrations", zap.Error(err))
	}

	repos := repositories.NewRepositories(db)
	jwtService := services.NewJWTService()
	authService := services.NewAuthService(repos.User, jwtService)

	defaultEmail := os.Getenv("DEFAULT_ADMIN_EMAIL")
	if defaultEmail == "" {
		defaultEmail = "admin@thomas.com"
	}

	defaultPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	if defaultPassword == "" {
		defaultPassword = "admin123"
	}

	if err := authService.CreateDefaultAdmin(defaultEmail, defaultPassword); err != nil {
		logger.Log.Error("Failed to create default admin", zap.Error(err))
	}

	// Create and initialize server
	server := servers.NewServer(db)
	server.InitializeRoutes()

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Log.Info("Server started successfully")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutdown signal received")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Server exited")
}
