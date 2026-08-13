package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDefault(t *testing.T) {
	var buf_1, buf_2 bytes.Buffer
	l1 := New(LevelInfo, &buf_1)
	l2 := New(LevelWarn, &buf_2)

	old := defaultLogger
	t.Cleanup(
		func() {
			defaultLogger = old
		},
	)

	SetDefault(l1)
	Info("hello %s", "world")
	assert.Contains(t, buf_1.String(), "hello world")
	assert.True(t, Enabled(LevelInfo))

	SetDefault(l2)
	Warn("second")
	assert.Contains(t, buf_2.String(), "second")
	assert.NotContains(t, buf_1.String(), "second")
	assert.False(t, Enabled(LevelInfo))
}

func TestPackageLevelFunctions(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelTrace, &buf)

	old := defaultLogger
	t.Cleanup(
		func() {
			defaultLogger = old
		},
	)
	SetDefault(l)

	tests := []struct {
		name   string
		logger func(string, ...any)
		label  string
	}{
		{
			name:   "Trace",
			logger: Trace,
			label:  "DEBUG",
		},
		{
			name:   "Debug",
			logger: Debug,
			label:  "DEBUG",
		},
		{
			name:   "Info",
			logger: Info,
			label:  "INFO",
		},
		{
			name:   "Warn",
			logger: Warn,
			label:  "WARN",
		},
		{
			name:   "Error",
			logger: Error,
			label:  "ERROR",
		},
		{
			name: "Log",
			logger: func(format string, args ...any) {
				Log(LevelInfo, format, args...)
			},
			label: "INFO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logger("hello %s", "world")
			output := buf.String()
			assert.Contains(t, output, tt.label)
			assert.Contains(t, output, "hello world")
		})
	}
}

func TestExit(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelError, &buf)

	oldLogger := defaultLogger
	t.Cleanup(
		func() {
			defaultLogger = oldLogger
		},
	)
	SetDefault(l)

	var exitCode int
	oldExit := exit
	exit = func(code int) {
		exitCode = code
	}
	t.Cleanup(
		func() {
			exit = oldExit
		},
	)

	Exit("fatal %s", "error")
	assert.Equal(t, 1, exitCode)
	output := buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, "fatal error")
}
