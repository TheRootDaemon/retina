package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultLoggerConfig(t *testing.T) {
	cfg := DefaultLoggerConfig()

	assert.Equal(t, "info", cfg.Level)
	assert.NotEmpty(t, cfg.File)
	assert.False(t, cfg.Enabled)
}

func TestDefaultLogFile(t *testing.T) {
	tests := []struct {
		name       string
		envXDG     string
		wantSuffix string
		wantAbs    bool
	}{
		{
			name:       "uses XDG_STATE_HOME when set",
			envXDG:     "STATE_DIR",
			wantSuffix: filepath.Join("hop", "hop.log"),
			wantAbs:    true,
		},
		{
			name:       "falls back to home dir",
			envXDG:     "",
			wantSuffix: filepath.Join("hop", "hop.log"),
			wantAbs:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envXDG != "" {
				dir := t.TempDir()
				t.Setenv("XDG_STATE_HOME", dir)
			} else {
				t.Setenv("XDG_STATE_HOME", "")
			}

			file, err := DefaultLogFile()
			require.NoError(t, err)
			assert.NotEmpty(t, file)
			if tt.wantAbs {
				assert.True(t, filepath.IsAbs(file))
			}
			assert.True(t, filepath.IsAbs(file))
			assert.Contains(t, file, "hop.log")
		})
	}
}

func TestFallbackLogFile(t *testing.T) {
	tests := []struct {
		name    string
		envHome string
		wantErr bool
	}{
		{
			name:    "valid home dir",
			envHome: "",
			wantErr: false,
		},
		{
			name:    "empty home dir",
			envHome: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "empty home dir" {
				t.Setenv("HOME", "")
				_, err := FallbackLogFile()
				require.Error(t, err)
				return
			}

			file, err := FallbackLogFile()
			require.NoError(t, err)
			assert.NotEmpty(t, file)
			assert.True(t, filepath.IsAbs(file))
			assert.Contains(t, file, "hop.log")
			assert.Contains(t, file, ".hop")
		})
	}
}

func TestDefaultLogFile_UsesXDGStateHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)

	file, err := DefaultLogFile()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "hop", "hop.log"), file)
}

func TestLoggerConfig_TOMLRoundTrip(t *testing.T) {
	input := `[logger]
enabled = true
level = "trace"
file = "/tmp/custom.log"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte(input), 0o600)
	require.NoError(t, err)

	cfg, err := LoadConfig(path)
	require.NoError(t, err)

	assert.True(t, cfg.Logger.Enabled)
	assert.Equal(t, "trace", cfg.Logger.Level)
	assert.Equal(t, "/tmp/custom.log", cfg.Logger.File)
}
