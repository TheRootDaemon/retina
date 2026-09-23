package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestZapLevel(t *testing.T) {
	tests := []struct {
		level Level
		want  zapcore.Level
	}{
		{
			level: LevelTrace,
			want:  zapcore.DebugLevel,
		},
		{
			level: LevelDebug,
			want:  zapcore.DebugLevel,
		},
		{
			level: LevelInfo,
			want:  zapcore.InfoLevel,
		},
		{
			level: LevelWarn,
			want:  zapcore.WarnLevel,
		},
		{
			level: LevelError,
			want:  zapcore.ErrorLevel,
		},
		{
			level: Level(999),
			want:  zapcore.InfoLevel,
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, zapLevel(tt.level))
	}
}

func TestLoggerClose(t *testing.T) {
	l := New(LevelInfo, new(bytes.Buffer))
	assert.NoError(t, l.Close())
}

func TestLoggerEnabled(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		arg   Level
		want  bool
	}{
		{
			name:  "error when trace",
			level: LevelError,
			arg:   LevelTrace,
			want:  false,
		},
		{
			name:  "error when debug",
			level: LevelError,
			arg:   LevelDebug,
			want:  false,
		},
		{
			name:  "error when info",
			level: LevelError,
			arg:   LevelInfo,
			want:  false,
		},
		{
			name:  "error when warn",
			level: LevelError,
			arg:   LevelWarn,
			want:  false,
		},
		{
			name:  "error when error",
			level: LevelError,
			arg:   LevelError,
			want:  true,
		},
		{
			name:  "info when trace",
			level: LevelInfo,
			arg:   LevelTrace,
			want:  false,
		},
		{
			name:  "info when debug",
			level: LevelInfo,
			arg:   LevelDebug,
			want:  false,
		},
		{
			name:  "info when info",
			level: LevelInfo,
			arg:   LevelInfo,
			want:  true,
		},
		{
			name:  "info when warn",
			level: LevelInfo,
			arg:   LevelWarn,
			want:  true,
		},
		{
			name:  "info when error",
			level: LevelInfo,
			arg:   LevelError,
			want:  true,
		},
		{
			name:  "trace when trace",
			level: LevelTrace,
			arg:   LevelTrace,
			want:  true,
		},
		{
			name:  "trace when debug",
			level: LevelTrace,
			arg:   LevelDebug,
			want:  true,
		},
		{
			name:  "trace when info",
			level: LevelTrace,
			arg:   LevelInfo,
			want:  true,
		},
		{
			name:  "warn when warn",
			level: LevelWarn,
			arg:   LevelWarn,
			want:  true,
		},
		{
			name:  "warn when trace",
			level: LevelWarn,
			arg:   LevelTrace,
			want:  false,
		},
		{
			name:  "warn when error",
			level: LevelWarn,
			arg:   LevelError,
			want:  true,
		},
		{
			name:  "very low level",
			level: LevelError,
			arg:   Level(-100),
			want:  false,
		},
		{
			name:  "very high level",
			level: LevelTrace,
			arg:   Level(999),
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.level, new(bytes.Buffer))
			assert.Equal(t, tt.want, l.Enabled(tt.arg))
		})
	}
}

func TestLogger_Emits(t *testing.T) {
	tests := []struct {
		name  string
		level string
		call  func(*Logger, string)
	}{
		{
			name:  "trace message",
			level: "trace",
			call:  func(l *Logger, s string) { l.Trace(s) },
		},
		{
			name:  "debug message",
			level: "debug",
			call:  func(l *Logger, s string) { l.Debug(s) },
		},
		{
			name:  "informational message",
			level: "info",
			call:  func(l *Logger, s string) { l.Info(s) },
		},
		{
			name:  "warning message",
			level: "warn",
			call:  func(l *Logger, s string) { l.Warn(s) },
		},
		{
			name:  "error message",
			level: "error",
			call:  func(l *Logger, s string) { l.Error(s) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := New(LevelTrace, &buf)
			tt.call(l, "hello"+tt.name)

			got := strings.ToLower(buf.String())
			assert.Contains(t, got, tt.level)
			assert.Contains(t, got, "hello")
		})
	}
}

func TestLogger_RespectsLevelledLogging(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelInfo, &buf)

	l.Trace("trace msg")
	l.Debug("debug msg")
	l.Info("info msg")
	l.Warn("warn msg")
	l.Error("error msg")

	out := buf.String()
	assert.NotContains(t, out, "trace")
	assert.NotContains(t, out, "debug")
	assert.Contains(t, out, "info")
	assert.Contains(t, out, "warn")
	assert.Contains(t, out, "error")
}

func TestLogger_SupportsFormatting(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelInfo, &buf)

	l.Info("hello %s", "world")
	assert.Contains(t, buf.String(), "hello world")
}
