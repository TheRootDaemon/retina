package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDefault(t *testing.T) {
	oldLogger := defaultLogger
	t.Cleanup(func() {
		defaultLogger = oldLogger
	})

	var buf_1, buf_2 bytes.Buffer
	logger_1 := New(LevelInfo, &buf_1)
	logger_2 := New(LevelWarn, &buf_2)

	SetDefault(logger_1)
	Info("log from logger_1")
	assert.True(t, Enabled(LevelInfo))
	assert.Contains(t, buf_1.String(), "logger_1")

	SetDefault(logger_2)
	Warn("log from logger_2")
	assert.False(t, Enabled(LevelInfo))
	assert.Contains(t, buf_2.String(), "logger_2")
	assert.NotContains(t, buf_1.String(), "logger_2")
}

func TestLogger(t *testing.T) {
	oldLogger := defaultLogger
	t.Cleanup(func() {
		defaultLogger = oldLogger
		assert.NoError(t, Close())
	})

	var buf bytes.Buffer
	SetDefault(New(LevelInfo, &buf))

	t.Run("emits enabled levels", func(t *testing.T) {
		Info("informational message")
		Warn("warning message")
		Error("error message")

		out := strings.ToLower(buf.String())
		for _, message := range []string{"informational message", "warning message", "error message"} {
			assert.Contains(t, out, message)
		}
	})

	t.Run("omits disabled levels", func(t *testing.T) {
		Trace("trace message")
		Debug("debug message")

		out := strings.ToLower(buf.String())
		for _, message := range []string{"trace message", "debug message"} {
			assert.NotContains(t, out, message)
		}
	})

	t.Run("supports formats", func(t *testing.T) {
		Info("hello %s", "world")

		out := strings.ToLower(buf.String())
		assert.Contains(t, out, "hello world")
	})

	t.Run("enabled levels", func(t *testing.T) {
		assert.False(t, Enabled(LevelTrace))
		assert.False(t, Enabled(LevelDebug))
		assert.True(t, Enabled(LevelInfo))
		assert.True(t, Enabled(LevelWarn))
		assert.True(t, Enabled(LevelError))
	})
}

func TestExit(t *testing.T) {
	oldLogger := defaultLogger
	oldExit := exit
	t.Cleanup(func() {
		defaultLogger = oldLogger
		exit = oldExit
	})

	var buf bytes.Buffer
	l := New(LevelError, &buf)
	SetDefault(l)

	var exitCode *int
	exit = func(code int) {
		exitCode = &code
	}

	Exit(11, "fatal %s", "error")
	assert.NotNil(t, exitCode)
	assert.Equal(t, 11, *exitCode)

	output := strings.ToLower(buf.String())
	assert.Contains(t, output, "error")
	assert.Contains(t, output, "fatal error")
}
