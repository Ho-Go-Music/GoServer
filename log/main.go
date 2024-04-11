package diylog

import (
	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	return sugar
}
