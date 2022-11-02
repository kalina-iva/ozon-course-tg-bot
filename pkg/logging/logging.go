package logging

import (
	"log"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger() error {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		return errors.Wrap(err, "cannot init logging")
	}
	zap.ReplaceGlobals(logger)
	return nil
}

func Close() {
	if err := logger.Sync(); err != nil {
		log.Print("sync logger failed")
	}
}
