package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level specifies the security level of a log.
type Level int

const (
	LevelTrace Level = iota - 2
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

// zapLevel maps a Level to a zapcore.Level.
func zapLevel(l Level) zapcore.Level {
	switch l {
	case LevelTrace:
		return zapcore.DebugLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Logger writes leveled log messages backed by uber-go/zap.
type Logger struct {
	sugar *zap.SugaredLogger
	level Level
}

// Close flushes any buffered log entries.
func (l *Logger) Close() error {
	return l.sugar.Sync()
}

// Enabled reports whether messages at the given level would be logged.
func (l *Logger) Enabled(level Level) bool {
	return l.level <= level
}

// Trace logs a message at LevelTrace.
func (l *Logger) Trace(format string, args ...any) {
	if !l.Enabled(LevelTrace) {
		return
	}
	l.sugar.Debugf(format, args...)
}

// Debug logs a message at LevelDebug.
func (l *Logger) Debug(format string, args ...any) {
	l.sugar.Debugf(format, args...)
}

// Info logs a message at LevelInfo.
func (l *Logger) Info(format string, args ...any) {
	l.sugar.Infof(format, args...)
}

// Warn logs a message at LevelWarn.
func (l *Logger) Warn(format string, args ...any) {
	l.sugar.Warnf(format, args...)
}

// Error logs a message at LevelError.
func (l *Logger) Error(format string, args ...any) {
	l.sugar.Errorf(format, args...)
}
