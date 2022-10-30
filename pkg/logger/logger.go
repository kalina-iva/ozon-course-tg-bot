package logger

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	env := "dev"
	if env == "dev" {
		logger, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		logger, err = cfg.Build()
	}
	if err != nil {
		log.Fatal("cannot init zap", err)
	}
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}
