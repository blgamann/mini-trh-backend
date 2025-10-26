package schemas

import (
	"mini-trh-backend/pkg/domain/entities"

	"gorm.io/gorm"
)

func MigrateUser(db *gorm.DB) error {
	if err := db.AutoMigrate(&entities.User{}); err != nil {
		return err
	}

	return nil
}
