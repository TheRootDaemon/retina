package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is the top-level configuration structure.
//
// Each field maps to a TOML section
// in the config file.
type Config struct {
	Salt   string       `toml:"salt"`
	Vault  string       `toml:"vault"`
	Logger LoggerConfig `toml:"logger"`
}

// Default returns a Config populated with default values.
func Default() Config {
	return Config{
		Logger: DefaultLoggerConfig(),
	}
}

// DefaultConfig generates the default configuration
// as a TOML-formatted string.
func DefaultConfig() (string, error) {
	var buf strings.Builder
	err := toml.NewEncoder(&buf).Encode(Default())
	return buf.String(), err
}

// ConfigPath returns the platform-appropriate
// path to the config file.
//
// The path can be overridden
// by setting the HOP_CONFIG environment variable.
//
// Default paths by platform:
//   - Linux:   ~/.config/hop/config.toml
//   - macOS:   ~/Library/Application Support/hop/config.toml
//   - Windows: %AppData%/hop/config.toml
func ConfigPath() string {
	if p := os.Getenv("HOP_CONFIG"); p != "" {
		return p
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, "hop", "config.toml")
}

// LoadConfig reads a TOML config file
// and returns the parsed Config.
//
// Fields not present in the file
// retain their default values.
func LoadConfig(path string) (*Config, error) {
	cfg := Default()
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
