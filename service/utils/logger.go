package utils

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

var Logger *otelzap.Logger

func InitializeLogger(daemon bool) {
	var (
		log *zap.Logger
		err error
	)
	if daemon {
		log, err = zap.NewProduction()
	} else {
		log, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}

	Logger = otelzap.New(log)
}
