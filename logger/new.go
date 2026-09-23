package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// newSugaredLogger builds a *zap.SugaredLogger that writes to w
// at the given minimum level using a human-readable console encoder.
func newSugaredLogger(w io.Writer, level Level) *zap.SugaredLogger {
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(w),
		zapLevel(level),
	)

	return zap.New(core).Sugar()
}

// New creates a Logger that writes to w at the given level.
func New(level Level, w io.Writer) *Logger {
	return &Logger{
		sugar: newSugaredLogger(w, level),
		level: level,
	}
}
