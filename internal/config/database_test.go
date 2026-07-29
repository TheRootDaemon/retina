package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultDatabaseConfig(t *testing.T) {
	tests := []struct {
		name   string
		envXDG string
	}{
		{
			name:   "uses XDG_STATE_HOME when set",
			envXDG: "STATE_DIR",
		},
		{
			name:   "XDG_STATE_HOME empty falls back to default",
			envXDG: "",
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

			cfg := DefaultDatabaseConfig()
			assert.NotEmpty(t, cfg.Path)
			if tt.envXDG != "" {
				assert.True(t, filepath.IsAbs(cfg.Path))
			}
		})
	}
}

func TestDatabaseConfigTOMLRoundTrip(t *testing.T) {
	input := `[database]
path = "/tmp/custom.db"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte(input), 0o600)
	require.NoError(t, err)

	cfg, err := LoadConfig(path)
	require.NoError(t, err)
	assert.Equal(t, "/tmp/custom.db", cfg.Database.Path)
}

func TestDatabaseConfigTOMLDefaults(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)

	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte(""), 0o644)
	require.NoError(t, err)

	cfg, err := LoadConfig(path)
	require.NoError(t, err)
	assert.Equal(t, DefaultDatabaseConfig(), cfg.Database)
}
