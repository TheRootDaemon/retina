package logger

import (
	"os"

	"github.com/therootdaemon/hop/internal/config"
	"go.uber.org/zap"
)

// defaultLogger is the package-level Logger used by the
// convenience functions. It starts as a no-op logger that
// discards output. Call SetDefault to replace it with a
// real logger after loading configuration.
var defaultLogger = &Logger{
	sugar: zap.NewNop().Sugar(),
	level: LevelInfo,
}

// SetDefault sets the package-level default logger.
func SetDefault(l *Logger) {
	defaultLogger = l
}

// SetDefaultFromConfig creates a Logger from config and
// sets it as the package-level default logger.
func SetDefaultFromConfig(cfg config.LoggerConfig) {
	SetDefault(NewFromConfig(cfg))
}

// Enabled reports whether the given level is enabled
// on the default logger.
func Enabled(level Level) bool {
	return defaultLogger.Enabled(level)
}

// Trace logs at LevelTrace via the default logger.
func Trace(format string, args ...any) {
	defaultLogger.Trace(format, args...)
}

// Debug logs at LevelDebug via the default logger.
func Debug(format string, args ...any) {
	defaultLogger.Debug(format, args...)
}

// Info logs at LevelInfo via the default logger.
func Info(format string, args ...any) {
	defaultLogger.Info(format, args...)
}

// Warn logs at LevelWarn via the default logger.
func Warn(format string, args ...any) {
	defaultLogger.Warn(format, args...)
}

// Error logs at LevelError via the default logger.
func Error(format string, args ...any) {
	defaultLogger.Error(format, args...)
}

// Exit logs at LevelError via the default logger and
// terminates the process with code 1.
func Exit(format string, args ...any) {
	defaultLogger.Error(format, args...)
	exit(1)
}

// Close flushes any buffered log entries on the default logger.
func Close() error {
	return defaultLogger.Close()
}

// Log logs at the given level via the default logger.
func Log(level Level, format string, args ...any) {
	defaultLogger.Log(level, format, args...)
}

// exit is a package-level variable so tests can replace it.
var exit = func(code int) {
	_ = defaultLogger.Close()
	os.Exit(code)
}
