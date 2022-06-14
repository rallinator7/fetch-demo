package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "ts",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.EpochTimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

// New creates a zap logger with a custom encoder config
func New() (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig = encoderConfig

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("NewLogger: %w", err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	return sugar, nil
}
