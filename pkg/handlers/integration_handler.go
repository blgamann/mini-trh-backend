package handlers

import (
	"mini-trh-backend/internal/utils"
	"mini-trh-backend/pkg/api/dto"
	"mini-trh-backend/pkg/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IntegrationHandler struct {
	integrationService *services.IntegrationService
}

func NewIntegrationHandler(integrationService *services.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{
		integrationService: integrationService,
	}
}

// InstallIntegration installs an integration for a node
// @Summary Install integration
// @Tags integrations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Node ID"
// @Param request body dto.InstallIntegrationRequest true "Integration data"
// @Success 201 {object} dto.IntegrationResponse
// @Router /nodes/{id}/integrations [post]
func (h *IntegrationHandler) InstallIntegration(c *gin.Context) {
	nodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid node ID")
		return
	}

	var req dto.InstallIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	integration, err := h.integrationService.InstallIntegration(nodeID, req.Type, req.Config)
	if err != nil {
		if err == services.ErrNodeNotFound {
			utils.NotFoundResponse(c, "Node not found")
			return
		}
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	response := dto.IntegrationResponse{
		ID:        integration.ID,
		NodeID:    integration.NodeID,
		Type:      integration.Type,
		Status:    integration.Status,
		Config:    integration.Config,
		Info:      integration.Info,
		LogPath:   integration.LogPath,
		CreatedAt: integration.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: integration.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	utils.CreatedResponse(c, response)
}

// UninstallIntegration removes an integration
// @Summary Uninstall integration
// @Tags integrations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Node ID"
// @Param type query string true "Integration type"
// @Success 200 {object} map[string]string
// @Router /nodes/{id}/integrations [delete]
func (h *IntegrationHandler) UninstallIntegration(c *gin.Context) {
	nodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid node ID")
		return
	}

	integrationType := c.Query("type")
	if integrationType == "" {
		utils.BadRequestResponse(c, "Integration type required")
		return
	}

	if err := h.integrationService.UninstallIntegration(nodeID, integrationType); err != nil {
		if err == services.ErrIntegrationNotFound {
			utils.NotFoundResponse(c, "Integration not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to uninstall integration")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Integration removal initiated"})
}

// GetIntegration retrieves an integration by ID
// @Summary Get integration by ID
// @Tags integrations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Node ID"
// @Param integrationId path string true "Integration ID"
// @Success 200 {object} dto.IntegrationResponse
// @Router /nodes/{id}/integrations/{integrationId} [get]
func (h *IntegrationHandler) GetIntegration(c *gin.Context) {
	integrationID, err := uuid.Parse(c.Param("integrationId"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid integration ID")
		return
	}

	integration, err := h.integrationService.GetIntegration(integrationID)
	if err != nil {
		if err == services.ErrIntegrationNotFound {
			utils.NotFoundResponse(c, "Integration not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to get integration")
		return
	}

	response := dto.IntegrationResponse{
		ID:        integration.ID,
		NodeID:    integration.NodeID,
		Type:      integration.Type,
		Status:    integration.Status,
		Config:    integration.Config,
		Info:      integration.Info,
		LogPath:   integration.LogPath,
		CreatedAt: integration.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: integration.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	utils.SuccessResponse(c, response)
}

// ListIntegrations retrieves all integrations for a node
// @Summary List integrations for a node
// @Tags integrations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Node ID"
// @Success 200 {object} dto.ListIntegrationsResponse
// @Router /nodes/{id}/integrations [get]
func (h *IntegrationHandler) ListIntegrations(c *gin.Context) {
	nodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid node ID")
		return
	}

	integrations, err := h.integrationService.ListIntegrations(nodeID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to list integrations")
		return
	}

	var integrationResponses []dto.IntegrationResponse
	for _, i := range integrations {
		integrationResponses = append(integrationResponses, dto.IntegrationResponse{
			ID:        i.ID,
			NodeID:    i.NodeID,
			Type:      i.Type,
			Status:    i.Status,
			Config:    i.Config,
			Info:      i.Info,
			LogPath:   i.LogPath,
			CreatedAt: i.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: i.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	response := dto.ListIntegrationsResponse{
		Integrations: integrationResponses,
		Total:        len(integrationResponses),
	}

	utils.SuccessResponse(c, response)
}
