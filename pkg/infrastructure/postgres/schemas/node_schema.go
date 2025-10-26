package schemas

import (
	"mini-trh-backend/pkg/domain/entities"

	"gorm.io/gorm"
)

// MigrateNode creates the nodes table and indexes
func MigrateNode(db *gorm.DB) error {
	// Auto-migrate
	if err := db.AutoMigrate(&entities.Node{}); err != nil {
		return err
	}

	// Additional indexes
	// Status index is already created via struct tag
	// Network index can be useful for filtering
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_nodes_network ON nodes(network)").Error; err != nil {
		return err
	}

	return nil
}
