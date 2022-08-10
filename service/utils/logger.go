package utils

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

var Logger *otelzap.Logger

func InitializeLogger() {
	var (
		log *zap.Logger
		err error
	)

	log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	Logger = otelzap.New(log)
}
