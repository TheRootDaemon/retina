package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therootdaemon/hop/internal/config"
	"go.uber.org/zap/zapcore"
)

func TestZapLevel(t *testing.T) {
	tests := []struct {
		level Level
		want  zapcore.Level
	}{
		{level: LevelTrace, want: zapcore.DebugLevel},
		{level: LevelDebug, want: zapcore.DebugLevel},
		{level: LevelInfo, want: zapcore.InfoLevel},
		{level: LevelWarn, want: zapcore.WarnLevel},
		{level: LevelError, want: zapcore.ErrorLevel},
		{level: Level(999), want: zapcore.InfoLevel},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, zapLevel(tt.level))
	}
}

func TestNewFromConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfg       config.LoggerConfig
		wantLevel Level
		wantNop   bool
		badState  bool
	}{
		{
			name:      "valid config writes to file",
			cfg:       config.LoggerConfig{Enabled: true, Level: "info", File: "hop/hop.log"},
			wantLevel: LevelInfo,
		},
		{
			name:      "disabled returns nop",
			cfg:       config.LoggerConfig{Enabled: false, Level: "debug", File: "hop/hop.log"},
			wantLevel: LevelDebug,
			wantNop:   true,
		},
		{
			name:      "empty file returns nop",
			cfg:       config.LoggerConfig{Enabled: true, Level: "warn", File: ""},
			wantLevel: LevelWarn,
			wantNop:   true,
		},
		{
			name:      "invalid level falls back to info",
			cfg:       config.LoggerConfig{Enabled: false, Level: "bogus", File: ""},
			wantLevel: LevelInfo,
			wantNop:   true,
		},
		{
			name:      "bad state dir falls back to nop",
			cfg:       config.LoggerConfig{Enabled: true, Level: "info", File: "hop/hop.log"},
			wantLevel: LevelInfo,
			wantNop:   true,
			badState:  true,
		},
		{
			name:      "disabled with empty file",
			cfg:       config.LoggerConfig{Enabled: false, Level: "trace", File: ""},
			wantLevel: LevelTrace,
			wantNop:   true,
		},
		{
			name:      "enabled with all valid fields",
			cfg:       config.LoggerConfig{Enabled: true, Level: "error", File: "hop/test.log"},
			wantLevel: LevelError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.badState {
				t.Setenv("XDG_STATE_HOME", "")
				t.Setenv("HOME", "/nonexistent")
			} else {
				t.Setenv("XDG_STATE_HOME", dir)
			}

			l := NewFromConfig(tt.cfg)
			require.NotNil(t, l)
			assert.Equal(t, tt.wantLevel, l.level)

			if tt.wantNop {
				l.Trace("nop %s", "trace")
				l.Debug("nop %s", "debug")
				l.Info("nop %s", "info")
				l.Warn("nop %s", "warn")
				l.Error("nop %s", "error")
				l.Log(LevelInfo, "nop %s", "log")
			} else {
				l.Error("write-check %d", 42)
			}
			require.NoError(t, l.Close())

			if !tt.wantNop {
				content, err := os.ReadFile(filepath.Join(dir, tt.cfg.File))
				require.NoError(t, err)
				assert.Contains(t, string(content), "write-check 42")
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input   string
		want    Level
		wantErr bool
	}{
		{input: "trace", want: LevelTrace},
		{input: "debug", want: LevelDebug},
		{input: "info", want: LevelInfo},
		{input: "warn", want: LevelWarn},
		{input: "error", want: LevelError},
		{input: "TRACE", want: LevelTrace},
		{input: "DEBUG", want: LevelDebug},
		{input: "INFO", want: LevelInfo},
		{input: "WARN", want: LevelWarn},
		{input: "ERROR", want: LevelError},
		{input: "Trace", want: LevelTrace},
		{input: "Debug", want: LevelDebug},
		{input: "Info", want: LevelInfo},
		{input: "Warn", want: LevelWarn},
		{input: "Error", want: LevelError},
		{input: "eRRoR", want: LevelError},
		{input: "MiXeD", want: LevelInfo, wantErr: true},
		{input: "unknown", want: LevelInfo, wantErr: true},
		{input: "", want: LevelInfo, wantErr: true},
		{input: "trace ", want: LevelInfo, wantErr: true},
		{input: " info", want: LevelInfo, wantErr: true},
		{input: "trace_debug", want: LevelInfo, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseLevel(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnabled(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		arg   Level
		want  bool
	}{
		{name: "error when trace", level: LevelError, arg: LevelTrace, want: false},
		{name: "error when debug", level: LevelError, arg: LevelDebug, want: false},
		{name: "error when info", level: LevelError, arg: LevelInfo, want: false},
		{name: "error when warn", level: LevelError, arg: LevelWarn, want: false},
		{name: "error when error", level: LevelError, arg: LevelError, want: true},
		{name: "info when trace", level: LevelInfo, arg: LevelTrace, want: false},
		{name: "info when debug", level: LevelInfo, arg: LevelDebug, want: false},
		{name: "info when info", level: LevelInfo, arg: LevelInfo, want: true},
		{name: "info when warn", level: LevelInfo, arg: LevelWarn, want: true},
		{name: "info when error", level: LevelInfo, arg: LevelError, want: true},
		{name: "trace when trace", level: LevelTrace, arg: LevelTrace, want: true},
		{name: "trace when debug", level: LevelTrace, arg: LevelDebug, want: true},
		{name: "trace when info", level: LevelTrace, arg: LevelInfo, want: true},
		{name: "warn when warn", level: LevelWarn, arg: LevelWarn, want: true},
		{name: "warn when trace", level: LevelWarn, arg: LevelTrace, want: false},
		{name: "warn when error", level: LevelWarn, arg: LevelError, want: true},
		{name: "very low level", level: LevelError, arg: Level(-100), want: false},
		{name: "very high level", level: LevelTrace, arg: Level(999), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.level, new(bytes.Buffer))
			assert.Equal(t, tt.want, l.Enabled(tt.arg))
		})
	}
}

func TestNewNop(t *testing.T) {
	tests := []struct {
		name  string
		level Level
	}{
		{name: "trace", level: LevelTrace},
		{name: "info", level: LevelInfo},
		{name: "error", level: LevelError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newNop(tt.level)
			require.NotNil(t, l)
			assert.Equal(t, tt.level, l.level)

			l.Trace("trace %s", "msg")
			l.Debug("debug %s", "msg")
			l.Info("info %s", "msg")
			l.Warn("warn %s", "msg")
			l.Error("error %s", "msg")
			l.Log(LevelInfo, "log %s", "msg")

			assert.NoError(t, l.Close())
		})
	}
}

func TestOpenLogFile(t *testing.T) {
	t.Run("creates file and dirs", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("XDG_STATE_HOME", dir)

		f, err := openLogFile("hop/hop.log")
		require.NoError(t, err)
		require.NotNil(t, f)
		_ = f.Close()

		info, err := os.Stat(filepath.Join(dir, "hop", "hop.log"))
		require.NoError(t, err)
		assert.True(t, info.Mode().IsRegular())
	})

	t.Run("missing state dir returns error", func(t *testing.T) {
		t.Setenv("XDG_STATE_HOME", "")
		t.Setenv("HOME", "/nonexistent")

		_, err := openLogFile("hop/hop.log")
		require.Error(t, err)
	})
}
