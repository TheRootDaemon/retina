package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	assert.NotNil(t, defaultLogger)
	assert.Equal(t, LevelInfo, defaultLogger.level)
}

func TestSetDefault(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelInfo, &buf)

	old := defaultLogger
	t.Cleanup(func() { defaultLogger = old })

	SetDefault(l)
	Info("hello %s", "world")

	output := buf.String()
	assert.Contains(t, output, "INFO")
	assert.Contains(t, output, "hello world")
}

func TestSetDefaultReplacesPrevious(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	l1 := New(LevelInfo, &buf1)
	l2 := New(LevelWarn, &buf2)

	old := defaultLogger
	t.Cleanup(
		func() {
			defaultLogger = old
		},
	)

	SetDefault(l1)
	Info("first")
	assert.Contains(t, buf1.String(), "first")

	SetDefault(l2)
	Warn("second")
	assert.Contains(t, buf2.String(), "second")

	// buf1 should NOT contain the second message (it goes to l2 now).
	assert.NotContains(t, buf1.String(), "second")
}

func TestPackageLevelFunctions(t *testing.T) {
	var buf bytes.Buffer
	l := New(LevelTrace, &buf)

	old := defaultLogger
	t.Cleanup(func() { defaultLogger = old })
	SetDefault(l)

	tests := []struct {
		name    string
		logFunc func(string, ...any)
		label   string
	}{
		{name: "Trace", logFunc: Trace, label: "DEBUG"},
		{name: "Debug", logFunc: Debug, label: "DEBUG"},
		{name: "Info", logFunc: Info, label: "INFO"},
		{name: "Warn", logFunc: Warn, label: "WARN"},
		{name: "Error", logFunc: Error, label: "ERROR"},
		{
			name: "Log",
			logFunc: func(format string, args ...any) {
				Log(LevelInfo, format, args...)
			},
			label: "INFO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc("hello %s", "world")
			output := buf.String()
			assert.Contains(t, output, tt.label)
			assert.Contains(t, output, "hello world")
		})
	}
}

func TestPackageLevelEnabled(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		arg      Level
		expected bool
	}{
		{name: "info when error", level: LevelError, arg: LevelInfo, expected: false},
		{name: "error when error", level: LevelError, arg: LevelError, expected: true},
		{name: "trace when info", level: LevelInfo, arg: LevelTrace, expected: false},
		{name: "info when info", level: LevelInfo, arg: LevelInfo, expected: true},
		{name: "warn when info", level: LevelInfo, arg: LevelWarn, expected: true},
		{name: "trace when trace", level: LevelTrace, arg: LevelDebug, expected: true},
		{name: "debug when warn", level: LevelWarn, arg: LevelDebug, expected: false},
		{name: "warn when warn", level: LevelWarn, arg: LevelWarn, expected: true},
		{name: "error when warn", level: LevelWarn, arg: LevelError, expected: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.level, new(bytes.Buffer))

			old := defaultLogger
			t.Cleanup(
				func() {
					defaultLogger = old
				},
			)
			SetDefault(logger)

			assert.Equal(t, tt.expected, Enabled(tt.arg))
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
		panic("exit")
	}
	t.Cleanup(
		func() {
			exit = oldExit
		},
	)

	assert.PanicsWithValue(
		t,
		"exit",
		func() {
			Exit("fatal %s", "error")
		},
	)

	assert.Equal(t, 1, exitCode)
	output := buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, "fatal error")
}
