package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogMethods(t *testing.T) {
	log := func(level Level) func(l *Logger, format string, args ...any) {
		return func(l *Logger, format string, args ...any) {
			l.Log(level, format, args...)
		}
	}

	tests := []struct {
		name      string
		minLevel  Level
		logger    func(l *Logger, format string, args ...any)
		wantLabel string
		wantWrite bool
	}{
		{
			name:      "trace",
			minLevel:  LevelTrace,
			logger:    (*Logger).Trace,
			wantLabel: "DEBUG",
			wantWrite: true,
		},
		{
			name:      "debug",
			minLevel:  LevelDebug,
			logger:    (*Logger).Debug,
			wantLabel: "DEBUG",
			wantWrite: true,
		},
		{
			name:      "info",
			minLevel:  LevelInfo,
			logger:    (*Logger).Info,
			wantLabel: "INFO",
			wantWrite: true,
		},
		{
			name:      "warn",
			minLevel:  LevelWarn,
			logger:    (*Logger).Warn,
			wantLabel: "WARN",
			wantWrite: true,
		},
		{
			name:      "error",
			minLevel:  LevelError,
			logger:    (*Logger).Error,
			wantLabel: "ERROR",
			wantWrite: true,
		},
		{
			name:      "log trace",
			minLevel:  LevelTrace,
			logger:    log(LevelTrace),
			wantLabel: "DEBUG",
			wantWrite: true,
		},
		{
			name:      "log debug",
			minLevel:  LevelDebug,
			logger:    log(LevelDebug),
			wantLabel: "DEBUG",
			wantWrite: true,
		},
		{
			name:      "log info",
			minLevel:  LevelInfo,
			logger:    log(LevelInfo),
			wantLabel: "INFO",
			wantWrite: true,
		},
		{
			name:      "log warn",
			minLevel:  LevelWarn,
			logger:    log(LevelWarn),
			wantLabel: "WARN",
			wantWrite: true,
		},
		{
			name:      "log error",
			minLevel:  LevelError,
			logger:    log(LevelError),
			wantLabel: "ERROR",
			wantWrite: true,
		},
		{
			name:     "trace disabled",
			minLevel: LevelDebug,
			logger:   (*Logger).Trace,
		},
		{
			name:     "debug disabled",
			minLevel: LevelInfo,
			logger:   (*Logger).Debug,
		},
		{
			name:     "info disabled",
			minLevel: LevelWarn,
			logger:   (*Logger).Info,
		},
		{
			name:     "warn disabled",
			minLevel: LevelError,
			logger:   (*Logger).Warn,
		},
		{
			name:     "log disabled",
			minLevel: LevelError,
			logger:   log(LevelInfo),
		},
		{
			name:     "log unknown level",
			minLevel: LevelTrace,
			logger:   log(Level(999)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := New(tt.minLevel, &buf)
			tt.logger(l, "hello %s", "world")

			if !tt.wantWrite {
				assert.Empty(t, buf.String())
				return
			}

			output := buf.String()
			assert.Contains(t, output, tt.wantLabel)
			assert.Contains(t, output, "hello world")
			assert.True(t, strings.HasSuffix(output, "\n"))
		})
	}
}
