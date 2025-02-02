package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mikaelchan/hamster/pkg/env"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func init() {
	level := zap.DebugLevel
	encoding := "console"
	development := true
	if env.IsRelease() {
		level = zap.InfoLevel
		encoding = "json"
		development = false
	}
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       development,
		Encoding:          encoding,
		DisableStacktrace: true,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ = config.Build()
	sugar = logger.Sugar()
}

func Debugf(format string, v ...any) {
	sugar.Debugf(format, v...)
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
