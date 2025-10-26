package postgres

import (
	"fmt"
	"os"
	"time"

	internalLogger "mini-trh-backend/internal/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Only log warnings and errors
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		internalLogger.Log.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		internalLogger.Log.Error("Failed to get SQL DB", zap.Error(err))
		return nil, err
	}

	// Connection pool settings (TRH Backend와 동일)
	sqlDB.SetMaxIdleConns(10)           // 유휴 연결 최대 개수
	sqlDB.SetMaxOpenConns(100)          // 최대 연결 개수
	sqlDB.SetConnMaxLifetime(time.Hour) // 연결 최대 수명

	internalLogger.Log.Info("Database connection established successfully")
	return db, nil
}
