package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/therootdaemon/hop/internal/config"
	"github.com/therootdaemon/hop/xdg"
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

// NewFromConfig creates a Logger from a LoggerConfig.
//
// If logging is disabled, the config has an empty file path,
// or the file cannot be opened, a no-op logger is returned
// that discards all output.
func NewFromConfig(cfg config.LoggerConfig) *Logger {
	level, err := ParseLevel(cfg.Level)
	if err != nil {
		level = LevelInfo
	}

	if !cfg.Enabled || cfg.File == "" {
		return newNop(level)
	}

	f, err := openLogFile(cfg.File)
	if err != nil {
		return newNop(level)
	}

	return New(level, f)
}

// ParseLevel parses a level name
// ("trace", "debug", "info", "warn", "error") into a Level.
// The comparison is case-insensitive.
// An unrecognized name returns LevelInfo
// and a non-nil error.
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(s) {
	case "trace":
		return LevelTrace, nil
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	default:
		return LevelInfo, fmt.Errorf("unknown log level: %q", s)
	}
}

// Close flushes any buffered log entries.
func (l *Logger) Close() error {
	return l.sugar.Sync()
}

// Enabled reports whether messages at the given level would be logged.
func (l *Logger) Enabled(level Level) bool {
	return l.level <= level
}

// newNop creates a disabled Logger that discards all output.
func newNop(level Level) *Logger {
	return &Logger{
		sugar: zap.NewNop().Sugar(),
		level: level,
	}
}

// openLogFile opens or creates a log file
// beneath the user's state directory.
//
// The provided path must be relative to the state directory
// returned by xdg.UserStateDir.
// Any missing parent directories are created automatically.
//
// The file is opened for writing in append mode
// and created if it does not already exists.
//
// The returned file has permissions 0600,
// and any newly created parent directories have permissions 0750.
//
// The state directory is opened with os.OpenRoot
// to prevent path traversal outside the user's state directory.
func openLogFile(path string) (*os.File, error) {
	stateDir, err := xdg.UserStateDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(
		filepath.Join(stateDir, filepath.Dir(path)),
		0o750,
	); err != nil {
		return nil, err
	}

	r, err := os.OpenRoot(stateDir)
	if err != nil {
		return nil, err
	}

	return r.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0o600,
	)
}
