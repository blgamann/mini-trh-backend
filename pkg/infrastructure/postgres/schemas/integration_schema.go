package schemas

import (
	"mini-trh-backend/pkg/domain/entities"

	"gorm.io/gorm"
)

// MigrateIntegration creates the integrations table and indexes
func MigrateIntegration(db *gorm.DB) error {
	if err := db.AutoMigrate(&entities.Integration{}); err != nil {
		return err
	}

	// Composite index for quick lookup by node and type
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_integrations_node_type ON integrations(node_id, type)").Error; err != nil {
		return err
	}

	return nil
}
