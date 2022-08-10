package utils

import "go.uber.org/zap"

var Logger *zap.Logger

func InitializeLogger() {
	var (
		err error
	)

	Logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}
