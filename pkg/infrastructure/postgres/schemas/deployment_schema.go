package schemas

import (
	"mini-trh-backend/pkg/domain/entities"

	"gorm.io/gorm"
)

// MigrateDeployment creates the deployments table and indexes
func MigrateDeployment(db *gorm.DB) error {
	if err := db.AutoMigrate(&entities.Deployment{}); err != nil {
		return err
	}

	// Index for filtering by step
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_deployments_step ON deployments(step)").Error; err != nil {
		return err
	}

	return nil
}
