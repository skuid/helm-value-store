/*
Package spec provides a zap.Logger creation function NewStandardLogger() for applications to use and write logs with.

NewStandardLogger() can be used like so:

    package main

    import "go.uber.org/zap"
    import "github.com/skuid/spec"

    func init() {
        l, err := spec.NewStandardLogger() // handle error
        zap.ReplaceGlobals(l)
    }

    func main() {
        zap.L().Debug("A debug message")
        zap.L().Info("An info message")
        zap.L().Info(
            "An info message with values",
            zap.String("key", "value"),
        )
        zap.L().Error("An error message")

        err := errors.New("some error")
        zap.L().Error("An error message", zap.Error(err))
    }

*/
package spec

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewStandardLogger creates a new zap.Logger based on common configuration
//
// This is intended to be used with zap.ReplaceGlobals() in an application's
// main.go.
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
