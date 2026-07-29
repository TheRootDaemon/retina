package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	assert.Equal(t, DefaultLoggerConfig(), cfg.Logger)
	assert.Empty(t, cfg.Salt)
	assert.Empty(t, cfg.Vault)
}

func TestDefaultConfig(t *testing.T) {
	s, err := DefaultConfig()
	require.NoError(t, err)
	require.NotEmpty(t, s)
}

func TestDefaultConfigRoundTrip(t *testing.T) {
	s, err := DefaultConfig()
	require.NoError(t, err)

	var cfg Config
	_, err = toml.Decode(s, &cfg)
	require.NoError(t, err)

	assert.Equal(t, DefaultLoggerConfig(), cfg.Logger)
	assert.Equal(t, Default().Salt, cfg.Salt)
	assert.Equal(t, Default().Vault, cfg.Vault)
}

func TestConfigPath(t *testing.T) {
	tests := []struct {
		name   string
		envVal string
		want   string
		wantFn func(t *testing.T) string
	}{
		{
			name:   "env var overrides default",
			envVal: "/custom/path/config.toml",
			want:   "/custom/path/config.toml",
		},
		{
			name:   "empty env var falls back to user config dir",
			envVal: "",
			wantFn: func(t *testing.T) string {
				dir, err := os.UserConfigDir()
				require.NoError(t, err)
				return filepath.Join(dir, "hop", "config.toml")
			},
		},
		{
			name:   "relative path from env var",
			envVal: "relative/config.toml",
			want:   "relative/config.toml",
		},
		{
			name:   "env var with spaces",
			envVal: "/path with spaces/config.toml",
			want:   "/path with spaces/config.toml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HOP_CONFIG", tt.envVal)
			if tt.wantFn != nil {
				assert.Equal(t, tt.wantFn(t), ConfigPath())
			} else {
				assert.Equal(t, tt.want, ConfigPath())
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, cfg *Config)
	}{
		{
			name:  "valid default",
			input: "",
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, DefaultLoggerConfig(), cfg.Logger)
				assert.Empty(t, cfg.Salt)
				assert.Empty(t, cfg.Vault)
			},
		},
		{
			name:  "salt and vault override",
			input: "salt = \"my-salt\"\nvault = \"/custom/vault\"\n",
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "my-salt", cfg.Salt)
				assert.Equal(t, "/custom/vault", cfg.Vault)
			},
		},
		{
			name:  "logger level override",
			input: "[logger]\nlevel = \"debug\"\n",
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "debug", cfg.Logger.Level)
				assert.Equal(t, DefaultLoggerConfig().File, cfg.Logger.File)
			},
		},
		{
			name:  "logger enabled override",
			input: "[logger]\nenabled = true\n",
			check: func(t *testing.T, cfg *Config) {
				assert.True(t, cfg.Logger.Enabled)
			},
		},
		{
			name:  "logger file override",
			input: "[logger]\nfile = \"/tmp/test.log\"\n",
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "/tmp/test.log", cfg.Logger.File)
			},
		},
		{
			name:    "malformed toml",
			input:   "invalid toml content {{{",
			wantErr: true,
		},
		{
			name:    "empty file",
			input:   "",
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, Default(), *cfg)
			},
		},
		{
			name:  "unknown top level keys ignored",
			input: "unknown_key = \"value\"\nanother = 42\n",
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, Default().Logger, cfg.Logger)
			},
		},
		{
			name:  "all fields overridden",
			input: "salt = \"s\"\nvault = \"v\"\n[logger]\nenabled = true\nlevel = \"warn\"\nfile = \"/custom.log\"\n",
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "s", cfg.Salt)
				assert.Equal(t, "v", cfg.Vault)
				assert.True(t, cfg.Logger.Enabled)
				assert.Equal(t, "warn", cfg.Logger.Level)
				assert.Equal(t, "/custom.log", cfg.Logger.File)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "config.toml")
			err := os.WriteFile(path, []byte(tt.input), 0o600)
			require.NoError(t, err)

			cfg, err := LoadConfig(path)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cfg)
				if tt.check != nil {
					tt.check(t, cfg)
				}
			}
		})
	}
}

func TestLoadConfig_fileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.toml")
	require.Error(t, err)
}

func TestLoadConfig_ReturnsNewPointer(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte("salt = \"a\"\n"), 0o600)
	require.NoError(t, err)

	cfg1, err := LoadConfig(path)
	require.NoError(t, err)
	cfg2, err := LoadConfig(path)
	require.NoError(t, err)

	assert.NotSame(t, cfg1, cfg2)
	assert.Equal(t, cfg1, cfg2)
}
