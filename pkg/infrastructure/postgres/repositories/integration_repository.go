package repositories

import (
	"errors"

	"mini-trh-backend/pkg/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IntegrationRepository handles database operations for integrations
type IntegrationRepository struct {
	db *gorm.DB
}

// NewIntegrationRepository creates a new IntegrationRepository instance
func NewIntegrationRepository(db *gorm.DB) *IntegrationRepository {
	return &IntegrationRepository{db: db}
}

// Create creates a new integration
func (r *IntegrationRepository) Create(integration *entities.Integration) error {
	return r.db.Create(integration).Error
}

// FindByID finds an integration by ID
func (r *IntegrationRepository) FindByID(id uuid.UUID) (*entities.Integration, error) {
	var integration entities.Integration
	err := r.db.Preload("Node").Where("id = ?", id).First(&integration).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &integration, nil
}

// FindByNodeID finds all integrations for a specific node
func (r *IntegrationRepository) FindByNodeID(nodeID uuid.UUID) ([]entities.Integration, error) {
	var integrations []entities.Integration
	err := r.db.Where("node_id = ?", nodeID).Find(&integrations).Error
	return integrations, err
}

// FindByNodeIDAndType finds a specific integration by node and type
func (r *IntegrationRepository) FindByNodeIDAndType(nodeID uuid.UUID, integrationType string) (*entities.Integration, error) {
	var integration entities.Integration
	err := r.db.Where("node_id = ? AND type = ?", nodeID, integrationType).First(&integration).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &integration, nil
}

// Update updates an integration
func (r *IntegrationRepository) Update(integration *entities.Integration) error {
	return r.db.Save(integration).Error
}

// UpdateStatus updates only the status field
func (r *IntegrationRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&entities.Integration{}).Where("id = ?", id).Update("status", status).Error
}

// Delete soft-deletes an integration
func (r *IntegrationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Integration{}, id).Error
}

// DeleteByNodeID deletes all integrations for a node
func (r *IntegrationRepository) DeleteByNodeID(nodeID uuid.UUID) error {
	return r.db.Where("node_id = ?", nodeID).Delete(&entities.Integration{}).Error
}

// ExistsByNodeIDAndType checks if an integration exists for a node
func (r *IntegrationRepository) ExistsByNodeIDAndType(nodeID uuid.UUID, integrationType string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Integration{}).
		Where("node_id = ? AND type = ?", nodeID, integrationType).
		Count(&count).Error
	return count > 0, err
}
