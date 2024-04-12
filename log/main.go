package diylog

import (
	"go.uber.org/zap"
	"sync"
)

var (
	Sugar *zap.SugaredLogger
	once  sync.Once
)

func init() {
	once.Do(func() {
		Sugar = newLogger()
	})
}
func newLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	return sugar
}
