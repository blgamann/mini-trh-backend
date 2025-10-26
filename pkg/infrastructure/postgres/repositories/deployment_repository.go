package repositories

import (
	"errors"

	"mini-trh-backend/pkg/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeploymentRepository handles database operations for deployments
type DeploymentRepository struct {
	db *gorm.DB
}

// NewDeploymentRepository creates a new DeploymentRepository instance
func NewDeploymentRepository(db *gorm.DB) *DeploymentRepository {
	return &DeploymentRepository{db: db}
}

// Create creates a new deployment record
func (r *DeploymentRepository) Create(deployment *entities.Deployment) error {
	return r.db.Create(deployment).Error
}

// FindByID finds a deployment by ID
func (r *DeploymentRepository) FindByID(id uuid.UUID) (*entities.Deployment, error) {
	var deployment entities.Deployment
	err := r.db.Preload("Node").Where("id = ?", id).First(&deployment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &deployment, nil
}

// FindByNodeID finds all deployments for a specific node
func (r *DeploymentRepository) FindByNodeID(nodeID uuid.UUID) ([]entities.Deployment, error) {
	var deployments []entities.Deployment
	err := r.db.Where("node_id = ?", nodeID).Order("created_at DESC").Find(&deployments).Error
	return deployments, err
}

// FindLatestByNodeID finds the most recent deployment for a node
func (r *DeploymentRepository) FindLatestByNodeID(nodeID uuid.UUID) (*entities.Deployment, error) {
	var deployment entities.Deployment
	err := r.db.Where("node_id = ?", nodeID).Order("created_at DESC").First(&deployment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &deployment, nil
}

// Update updates a deployment record
func (r *DeploymentRepository) Update(deployment *entities.Deployment) error {
	return r.db.Save(deployment).Error
}

// UpdateStatus updates only the status field
func (r *DeploymentRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&entities.Deployment{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateStatusAndError updates status and error message
func (r *DeploymentRepository) UpdateStatusAndError(id uuid.UUID, status, errorMsg string) error {
	return r.db.Model(&entities.Deployment{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"error_msg": errorMsg,
	}).Error
}

// Delete soft-deletes a deployment
func (r *DeploymentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Deployment{}, id).Error
}

// DeleteByNodeID deletes all deployments for a node
func (r *DeploymentRepository) DeleteByNodeID(nodeID uuid.UUID) error {
	return r.db.Where("node_id = ?", nodeID).Delete(&entities.Deployment{}).Error
}
