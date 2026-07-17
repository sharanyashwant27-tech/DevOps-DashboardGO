package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a production-ready Zap logger.
func New(mode string) (*zap.Logger, error) {
	var cfg zap.Config
	if mode == "release" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	return cfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// Sync flushes buffered logs.
func Sync(l *zap.Logger) {
	if l != nil {
		_ = l.Sync()
	}
}

// MustNew panics if logger cannot be created.
func MustNew(mode string) *zap.Logger {
	l, err := New(mode)
	if err != nil {
		panic(err)
	}
	return l
}

// WithPID adds process id field.
func WithPID(l *zap.Logger) *zap.Logger {
	return l.With(zap.Int("pid", os.Getpid()))
}
