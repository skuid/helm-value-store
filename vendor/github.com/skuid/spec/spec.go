/*
Package spec provides a Logger for applications to use and write logs with.

The Logger can be used like so:

	import "github.com/skuid/spec"
	import "github.com/uber-go/zap"

	var logger spec.Logger

	func main() {
		logger.Debug("A debug message")
		logger.Info("An info message")
		logger.Info(
			"An info message with values",
			zap.String("key", "value"),
		)
		logger.Error("An error message")

		err := errors.New("some error")
		logger.Error("An error message", zap.Error(err))
	}

*/
package spec

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Logger is a zap.Logger that can be set to the same logger you've created
var Logger *zap.Logger

//NewStandardLogger creates a new zap.Logger based on common configuration
func NewStandardLogger() (l *zap.Logger, err error) {
	config := zap.Config{
		Level:       zap.NewAtomicLevel(),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return config.Build()
}
