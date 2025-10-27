package servers

import (
	"context"
	"fmt"
	"mini-trh-backend/internal/logger"
	"mini-trh-backend/pkg/api/routes"
	"mini-trh-backend/pkg/handlers"
	"mini-trh-backend/pkg/infrastructure/postgres/repositories"
	"mini-trh-backend/pkg/middleware"
	"mini-trh-backend/pkg/services"
	"mini-trh-backend/pkg/services/taskmanager"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	Router      *gin.Engine
	DB          *gorm.DB
	TaskManager *taskmanager.TaskManager
	httpServer  *http.Server
}

func NewServer(db *gorm.DB) *Server {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Add recovery middleware (panic recovery)
	router.Use(gin.Recovery())

	// Add logging middleware
	router.Use(middleware.LoggingMiddleware())

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.MaxAge = 12 * time.Hour
	router.Use(cors.New(config))

	taskManager := taskmanager.NewTaskManager(5, 20) // 5 workers, buffer size 20

	return &Server{
		Router:      router,
		DB:          db,
		TaskManager: taskManager,
	}
}

func (s *Server) InitializeRoutes() {
	logger.Log.Info("Initializing routes...")

	// Initialize repositories
	repos := repositories.NewRepositories(s.DB)

	// Initialize services
	jwtService := services.NewJWTService()
	authService := services.NewAuthService(repos.User, jwtService)
	nodeService := services.NewNodeDeploymentService(
		repos.Node,
		repos.Deployment,
		repos.Integration,
		s.TaskManager,
	)
	integrationService := services.NewIntegrationService(
		repos.Node,
		repos.Integration,
		s.TaskManager,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	nodeHandler := handlers.NewNodeHandler(nodeService)
	integrationHandler := handlers.NewIntegrationHandler(integrationService)
	healthHandler := handlers.NewHealthHandler()

	// Setup routes
	routes.SetupRoutes(
		s.Router,
		jwtService,
		authHandler,
		nodeHandler,
		integrationHandler,
		healthHandler,
	)

	logger.Log.Info("Routes initialized successfully")
}

func (s *Server) Start() error {
	// Start task manager
	s.TaskManager.Start()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      s.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Log.Info("Starting server", zap.String("port", port))

	// Start server
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Error("Failed to start server", zap.Error(err))
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Log.Info("Shutting down server...")

	// Stop task manager
	s.TaskManager.Stop()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Log.Error("Server shutdown error", zap.Error(err))
		return err
	}

	logger.Log.Info("Server shutdown complete")
	return nil
}
