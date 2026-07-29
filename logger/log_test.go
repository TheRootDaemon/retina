package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogMethods(t *testing.T) {
	tests := []struct {
		name      string
		minLevel  Level
		logFunc   func(l *Logger, format string, args ...any)
		logLevel  Level
		wantLabel string
		wantText  string
	}{
		{
			name:      "Trace at LevelTrace",
			minLevel:  LevelTrace,
			logFunc:   (*Logger).Trace,
			logLevel:  LevelTrace,
			wantLabel: "DEBUG",
			wantText:  "hello world",
		},
		{
			name:      "Debug at LevelDebug",
			minLevel:  LevelDebug,
			logFunc:   (*Logger).Debug,
			logLevel:  LevelDebug,
			wantLabel: "DEBUG",
			wantText:  "hello world",
		},
		{
			name:      "Info at LevelInfo",
			minLevel:  LevelInfo,
			logFunc:   (*Logger).Info,
			logLevel:  LevelInfo,
			wantLabel: "INFO",
			wantText:  "hello world",
		},
		{
			name:      "Warn at LevelWarn",
			minLevel:  LevelWarn,
			logFunc:   (*Logger).Warn,
			logLevel:  LevelWarn,
			wantLabel: "WARN",
			wantText:  "hello world",
		},
		{
			name:      "Error at LevelError",
			minLevel:  LevelError,
			logFunc:   (*Logger).Error,
			logLevel:  LevelError,
			wantLabel: "ERROR",
			wantText:  "hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := New(tt.minLevel, &buf)

			tt.logFunc(l, "hello %s", "world")

			output := buf.String()
			assert.Contains(t, output, tt.wantLabel)
			assert.Contains(t, output, tt.wantText)
			assert.True(t, strings.HasSuffix(output, "\n"))
		})
	}
}

func TestLogMethodsDisabled(t *testing.T) {
	tests := []struct {
		name     string
		minLevel Level
		logFunc  func(l *Logger, format string, args ...any)
	}{
		{name: "Trace disabled at LevelDebug", minLevel: LevelDebug, logFunc: (*Logger).Trace},
		{name: "Debug disabled at LevelInfo", minLevel: LevelInfo, logFunc: (*Logger).Debug},
		{name: "Info disabled at LevelWarn", minLevel: LevelWarn, logFunc: (*Logger).Info},
		{name: "Warn disabled at LevelError", minLevel: LevelError, logFunc: (*Logger).Warn},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := New(tt.minLevel, &buf)
			tt.logFunc(l, "hidden %s", "msg")

			assert.Empty(t, buf.String())
		})
	}
}

func TestLogMethod(t *testing.T) {
	tests := []struct {
		name      string
		minLevel  Level
		logLevel  Level
		wantLabel string
		wantText  string
		wantWrite bool
	}{
		{name: "Trace", minLevel: LevelTrace, logLevel: LevelTrace, wantLabel: "DEBUG", wantText: "msg", wantWrite: true},
		{name: "Debug", minLevel: LevelDebug, logLevel: LevelDebug, wantLabel: "DEBUG", wantText: "msg", wantWrite: true},
		{name: "Info", minLevel: LevelInfo, logLevel: LevelInfo, wantLabel: "INFO", wantText: "msg", wantWrite: true},
		{name: "Warn", minLevel: LevelWarn, logLevel: LevelWarn, wantLabel: "WARN", wantText: "msg", wantWrite: true},
		{name: "Error", minLevel: LevelError, logLevel: LevelError, wantLabel: "ERROR", wantText: "msg", wantWrite: true},
		{name: "disabled", minLevel: LevelError, logLevel: LevelInfo, wantWrite: false},
		{name: "unknown level", minLevel: LevelTrace, logLevel: Level(999), wantWrite: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := New(tt.minLevel, &buf)
			l.Log(tt.logLevel, "hello %s", "msg")

			if tt.wantWrite {
				output := buf.String()
				assert.Contains(t, output, tt.wantLabel)
				assert.Contains(t, output, tt.wantText)
			} else {
				assert.Empty(t, buf.String())
			}
		})
	}
}

func TestLogMethodsNoArgs(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelTrace, &buf)

	l.Trace("simple")
	l.Debug("simple")
	l.Info("simple")
	l.Warn("simple")
	l.Error("simple")

	output := buf.String()
	assert.Equal(t, 5, strings.Count(output, "simple"))
}

func TestLogMethodsFormatArgs(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelTrace, &buf)

	l.Info("count=%d name=%s", 42, "test")

	output := buf.String()
	assert.Contains(t, output, "count=42 name=test")
}
