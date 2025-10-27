package dto

import (
	"mini-trh-backend/pkg/domain/entities"

	"github.com/google/uuid"
)

// InstallIntegrationRequest represents an integration installation request
type InstallIntegrationRequest struct {
	Type   string         `json:"type" binding:"required,oneof=monitoring backup explorer"`
	Config entities.JSONB `json:"config" binding:"required"`
}

// IntegrationResponse represents an integration response
type IntegrationResponse struct {
	ID        uuid.UUID      `json:"id"`
	NodeID    uuid.UUID      `json:"node_id"`
	Type      string         `json:"type"`
	Status    string         `json:"status"`
	Config    entities.JSONB `json:"config,omitempty"`
	Info      entities.JSONB `json:"info,omitempty"`
	LogPath   string         `json:"log_path,omitempty"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

// ListIntegrationsResponse represents a list of integrations
type ListIntegrationsResponse struct {
	Integrations []IntegrationResponse `json:"integrations"`
	Total        int                   `json:"total"`
}
