package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger returns a custom logger for Shipswift to be used with all the services. This
// guarantees we have consistent logging across the platform
func NewLogger(serviceName, level string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "fatal":
		config.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.DisableStacktrace = true

	l, err := config.Build()
	if err != nil {
		return nil, err
	}

	return l.Named(serviceName).Sugar(), nil
}
