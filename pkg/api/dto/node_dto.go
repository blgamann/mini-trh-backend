package dto

import (
	"mini-trh-backend/pkg/domain/entities"

	"github.com/google/uuid"
)

// CreateNodeRequest represents a node creation request
type CreateNodeRequest struct {
	Name    string         `json:"name" binding:"required,min=3,max=50"`
	Network string         `json:"network" binding:"required,oneof=mainnet testnet"`
	Config  entities.JSONB `json:"config" binding:"required"`
}

// NodeResponse represents a node response
type NodeResponse struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Network   string         `json:"network"`
	Status    string         `json:"status"`
	Config    entities.JSONB `json:"config,omitempty"`
	Metadata  entities.JSONB `json:"metadata,omitempty"`
	UserID    uuid.UUID      `json:"user_id"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

// NodeDetailResponse includes related data
type NodeDetailResponse struct {
	NodeResponse
	Deployments  []DeploymentResponse  `json:"deployments,omitempty"`
	Integrations []IntegrationResponse `json:"integrations,omitempty"`
}

// ListNodesResponse represents a list of nodes
type ListNodesResponse struct {
	Nodes []NodeResponse `json:"nodes"`
	Total int            `json:"total"`
}
