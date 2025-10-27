package dto

import (
	"mini-trh-backend/pkg/domain/entities"

	"github.com/google/uuid"
)

// DeploymentResponse represents a deployment response
type DeploymentResponse struct {
	ID        uuid.UUID      `json:"id"`
	NodeID    uuid.UUID      `json:"node_id"`
	Step      string         `json:"step"`
	Status    string         `json:"status"`
	Config    entities.JSONB `json:"config,omitempty"`
	LogPath   string         `json:"log_path,omitempty"`
	ErrorMsg  string         `json:"error_msg,omitempty"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

// ListDeploymentsResponse represents a list of deployments
type ListDeploymentsResponse struct {
	Deployments []DeploymentResponse `json:"deployments"`
	Total       int                  `json:"total"`
}
