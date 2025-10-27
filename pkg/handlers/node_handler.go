package handlers

import (
	"mini-trh-backend/internal/utils"
	"mini-trh-backend/pkg/api/dto"
	"mini-trh-backend/pkg/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NodeHandler struct {
	nodeService *services.NodeDeploymentService
}

func NewNodeHandler(nodeService *services.NodeDeploymentService) *NodeHandler {
	return &NodeHandler{
		nodeService: nodeService,
	}
}

// CreateNode creates a new node
// @Summary Create a new node
// @Tags nodes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateNodeRequest true "Node creation data"
// @Success 201 {object} dto.NodeResponse
// @Router /nodes [post]
func (h *NodeHandler) CreateNode(c *gin.Context) {
	var req dto.CreateNodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	// Create node
	node, err := h.nodeService.CreateNode(userID.(uuid.UUID), req.Name, req.Network, req.Config)
	if err != nil {
		if err == services.ErrNodeAlreadyExists {
			utils.BadRequestResponse(c, "Node with this name already exists")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to create node")
		return
	}

	response := dto.NodeResponse{
		ID:        node.ID,
		Name:      node.Name,
		Network:   node.Network,
		Status:    node.Status,
		Config:    node.Config,
		Metadata:  node.Metadata,
		UserID:    node.UserID,
		CreatedAt: node.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: node.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	utils.CreatedResponse(c, response)
}

// GetNode retrieves a node by ID
// @Summary Get node by ID
// @Tags nodes
// @Produce json
// @Security BearerAuth
// @Param id path string true "Node ID"
// @Success 200 {object} dto.NodeDetailResponse
// @Router /nodes/{id} [get]
func (h *NodeHandler) GetNode(c *gin.Context) {
	nodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid node ID")
		return
	}

	node, err := h.nodeService.GetNode(nodeID)
	if err != nil {
		if err == services.ErrNodeNotFound {
			utils.NotFoundResponse(c, "Node not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to get node")
		return
	}

	// Build response with related data
	response := dto.NodeDetailResponse{
		NodeResponse: dto.NodeResponse{
			ID:        node.ID,
			Name:      node.Name,
			Network:   node.Network,
			Status:    node.Status,
			Config:    node.Config,
			Metadata:  node.Metadata,
			UserID:    node.UserID,
			CreatedAt: node.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: node.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}

	// Add deployments
	if node.Deployments != nil {
		for _, d := range node.Deployments {
			response.Deployments = append(response.Deployments, dto.DeploymentResponse{
				ID:        d.ID,
				NodeID:    d.NodeID,
				Step:      d.Step,
				Status:    d.Status,
				Config:    d.Config,
				LogPath:   d.LogPath,
				ErrorMsg:  d.ErrorMsg,
				CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z"),
				UpdatedAt: d.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			})
		}
	}

	// Add integrations
	if node.Integrations != nil {
		for _, i := range node.Integrations {
			response.Integrations = append(response.Integrations, dto.IntegrationResponse{
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
	}

	utils.SuccessResponse(c, response)
}

// ListNodes retrieves all nodes with pagination
// @Summary List all nodes
// @Tags nodes
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.ListNodesResponse
// @Router /nodes [get]
func (h *NodeHandler) ListNodes(c *gin.Context) {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get nodes
	nodes, err := h.nodeService.ListNodes(limit, offset, nil)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to list nodes")
		return
	}

	// Build response
	var nodeResponses []dto.NodeResponse
	for _, node := range nodes {
		nodeResponses = append(nodeResponses, dto.NodeResponse{
			ID:        node.ID,
			Name:      node.Name,
			Network:   node.Network,
			Status:    node.Status,
			Config:    node.Config,
			Metadata:  node.Metadata,
			UserID:    node.UserID,
			CreatedAt: node.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: node.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	response := dto.ListNodesResponse{
		Nodes: nodeResponses,
		Total: len(nodeResponses),
	}

	utils.SuccessResponse(c, response)
}

// DeleteNode deletes a node
// @Summary Delete a node
// @Tags nodes
// @Produce json
// @Security BearerAuth
// @Param id path string true "Node ID"
// @Success 200 {object} map[string]string
// @Router /nodes/{id} [delete]
func (h *NodeHandler) DeleteNode(c *gin.Context) {
	nodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid node ID")
		return
	}

	if err := h.nodeService.DeleteNode(nodeID); err != nil {
		if err == services.ErrNodeNotFound {
			utils.NotFoundResponse(c, "Node not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to delete node")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Node deletion initiated"})
}

// GetDeployment retrieves a deployment by ID
// @Summary Get deployment by ID
// @Tags nodes
// @Produce json
// @Security BearerAuth
// @Param id path string true "Node ID"
// @Param deploymentId path string true "Deployment ID"
// @Success 200 {object} dto.DeploymentResponse
// @Router /nodes/{id}/deployments/{deploymentId} [get]
func (h *NodeHandler) GetDeployment(c *gin.Context) {
	deploymentID, err := uuid.Parse(c.Param("deploymentId"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid deployment ID")
		return
	}

	deployment, err := h.nodeService.GetDeployment(deploymentID)
	if err != nil {
		if err == services.ErrDeploymentNotFound {
			utils.NotFoundResponse(c, "Deployment not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to get deployment")
		return
	}

	response := dto.DeploymentResponse{
		ID:        deployment.ID,
		NodeID:    deployment.NodeID,
		Step:      deployment.Step,
		Status:    deployment.Status,
		Config:    deployment.Config,
		LogPath:   deployment.LogPath,
		ErrorMsg:  deployment.ErrorMsg,
		CreatedAt: deployment.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: deployment.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	utils.SuccessResponse(c, response)
}

// ListDeployments retrieves all deployments for a node
// @Summary List deployments for a node
// @Tags nodes
// @Produce json
// @Security BearerAuth
// @Param id path string true "Node ID"
// @Success 200 {object} dto.ListDeploymentsResponse
// @Router /nodes/{id}/deployments [get]
func (h *NodeHandler) ListDeployments(c *gin.Context) {
	nodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid node ID")
		return
	}

	deployments, err := h.nodeService.ListDeployments(nodeID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to list deployments")
		return
	}

	var deploymentResponses []dto.DeploymentResponse
	for _, d := range deployments {
		deploymentResponses = append(deploymentResponses, dto.DeploymentResponse{
			ID:        d.ID,
			NodeID:    d.NodeID,
			Step:      d.Step,
			Status:    d.Status,
			Config:    d.Config,
			LogPath:   d.LogPath,
			ErrorMsg:  d.ErrorMsg,
			CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: d.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	response := dto.ListDeploymentsResponse{
		Deployments: deploymentResponses,
		Total:       len(deploymentResponses),
	}

	utils.SuccessResponse(c, response)
}
