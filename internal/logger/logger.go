package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() error {
	config := zap.NewProductionConfig()

	// config log level (Debug, Info, Warn, Error)
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 	2025-10-26T20:13:45.234+0900  INFO  api/server.go:42  Server started  {"port": 8080}
	// 2025-10-26T20:13:45.237+0900  DEBUG db/connection.go:18  Connected to DB  {"dsn": "postgres://localhost:5432"}
	config.Encoding = "console"

	var err error
	Log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func Sync() {
	if Log != nil {
		// ignore error
		_ = Log.Sync()
	}
}
