package routes

import (
	"mini-trh-backend/pkg/handlers"
	"mini-trh-backend/pkg/middleware"
	"mini-trh-backend/pkg/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	jwtService *services.JWTService,
	authHandler *handlers.AuthHandler,
	nodeHandler *handlers.NodeHandler,
	integrationHandler *handlers.IntegrationHandler,
	healthHandler *handlers.HealthHandler,
) {
	v1 := router.Group("/api/v1")

	v1.GET("/health", healthHandler.HealthCheck)

	authRoutes := v1.Group("/auth")
	{
		// Public routes
		authRoutes.POST("/login", authHandler.Login)

		authProtected := authRoutes.Group("")
		authProtected.Use(middleware.JWTMiddleware(jwtService))
		{
			authProtected.GET("/profile", authHandler.GetProfile)

			// Admin only routes
			authAdmin := authProtected.Group("")
			authAdmin.Use(middleware.AdminOnly())
			{
				authAdmin.POST("/register", authHandler.Register)
			}
		}
	}

	// Node routes
	nodeRoutes := v1.Group("/nodes")
	nodeRoutes.Use(middleware.JWTMiddleware(jwtService)) // All node routes require authentication
	{
		// List and get nodes (any authenticated user)
		nodeRoutes.GET("", nodeHandler.ListNodes)
		nodeRoutes.GET("/:id", nodeHandler.GetNode)

		// Admin only routes
		nodeAdmin := nodeRoutes.Group("")
		nodeAdmin.Use(middleware.AdminOnly())
		{
			nodeAdmin.POST("", nodeHandler.CreateNode)
			nodeAdmin.DELETE("/:id", nodeHandler.DeleteNode)
		}

		// Deployment routes (authenticated users can view)
		nodeRoutes.GET("/:id/deployments", nodeHandler.ListDeployments)
		nodeRoutes.GET("/:id/deployments/:deploymentId", nodeHandler.GetDeployment)

		// Integration routes
		integrationRoutes := nodeRoutes.Group("/:id/integrations")
		{
			// View integrations (any authenticated user)
			integrationRoutes.GET("", integrationHandler.ListIntegrations)
			integrationRoutes.GET("/:integrationId", integrationHandler.GetIntegration)

			// Install/uninstall (admin only)
			integrationAdmin := integrationRoutes.Group("")
			integrationAdmin.Use(middleware.AdminOnly())
			{
				integrationAdmin.POST("", integrationHandler.InstallIntegration)
				integrationAdmin.DELETE("", integrationHandler.UninstallIntegration)
			}
		}
	}
}
