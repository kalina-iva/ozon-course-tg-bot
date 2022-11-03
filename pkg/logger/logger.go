package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger(env string) error {
	var err error
	if env == "dev" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return errors.Wrap(err, "cannot init logger")
	}
	zap.ReplaceGlobals(logger)
	return nil
}

func Close() error {
	return logger.Sync()
}

func Info(message string, fields ...zap.Field) {
	zap.L().Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	zap.L().Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	zap.L().Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	zap.L().Fatal(message, fields...)
}
