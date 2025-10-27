package handlers

import (
	"mini-trh-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck returns the health status of the service
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, gin.H{
		"status":  "healthy",
		"service": "trh-backend-thomas",
	})
}
