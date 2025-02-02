package logger

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func init() {
	logger, _ = zap.NewProduction()
	sugar = logger.Sugar()
}

func Infof(format string, v ...any) {
	sugar.Infof(format, v...)
}

func Errorf(format string, v ...any) {
	sugar.Errorf(format, v...)
}

func Fatalf(format string, v ...any) {
	sugar.Fatalf(format, v...)
}

func Panicf(format string, v ...any) {
	sugar.Panicf(format, v...)
}

func Sync() {
	logger.Sync()
}
