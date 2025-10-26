package schemas

import (
	"mini-trh-backend/internal/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	logger.Log.Info("Running database migrations...")

	migrations := []func(*gorm.DB) error{
		MigrateUser,
		MigrateNode,
		MigrateDeployment,
		MigrateIntegration,
	}

	for _, migrate := range migrations {
		if err := migrate(db); err != nil {
			logger.Log.Error("Migration failed", zap.Error(err))
			return err
		}
	}

	logger.Log.Info("Database migrations completed successfully")
	return nil
}
