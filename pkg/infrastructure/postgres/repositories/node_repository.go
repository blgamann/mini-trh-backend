package repositories

import (
	"errors"

	"mini-trh-backend/pkg/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NodeRepository handles database operations for nodes
type NodeRepository struct {
	db *gorm.DB
}

// NewNodeRepository creates a new NodeRepository instance
func NewNodeRepository(db *gorm.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

// Create creates a new node
func (r *NodeRepository) Create(node *entities.Node) error {
	return r.db.Create(node).Error
}

// FindByID finds a node by ID with optional preloading
func (r *NodeRepository) FindByID(id uuid.UUID, preload ...string) (*entities.Node, error) {
	var node entities.Node
	query := r.db

	// Preload associations if specified
	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Where("id = ?", id).First(&node).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &node, nil
}

// FindAll retrieves all nodes with pagination and optional filtering
func (r *NodeRepository) FindAll(limit, offset int, filters map[string]interface{}) ([]entities.Node, error) {
	var nodes []entities.Node
	query := r.db.Limit(limit).Offset(offset)

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	err := query.Preload("User").Find(&nodes).Error
	return nodes, err
}

// FindByUserID finds all nodes belonging to a user
func (r *NodeRepository) FindByUserID(userID uuid.UUID) ([]entities.Node, error) {
	var nodes []entities.Node
	err := r.db.Where("user_id = ?", userID).Preload("User").Find(&nodes).Error
	return nodes, err
}

// Update updates a node's information
func (r *NodeRepository) Update(node *entities.Node) error {
	return r.db.Save(node).Error
}

// UpdateStatus updates only the status field
func (r *NodeRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&entities.Node{}).Where("id = ?", id).Update("status", status).Error
}

// Delete soft-deletes a node
func (r *NodeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Node{}, id).Error
}

// Count returns the total number of nodes
func (r *NodeRepository) Count(filters map[string]interface{}) (int64, error) {
	query := r.db.Model(&entities.Node{})

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// ExistsByName checks if a node with the given name exists
func (r *NodeRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Node{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
