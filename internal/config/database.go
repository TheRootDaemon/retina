package config

import (
	"path/filepath"

	"github.com/therootdaemon/hop/xdg"
)

// DatabaseConfig holds the database connection configuration.
// It maps to the [database] section in the config.
type DatabaseConfig struct {
	// Path is the path to the SQLite database file.
	// By default it is located under the user's state directory
	// as "<state-dir>/hop/hop.db".
	Path string `toml:"path"`
}

// DefaultDatabaseConfig returns the default database settings.
func DefaultDatabaseConfig() DatabaseConfig {
	dir, err := xdg.UserStateDir()
	if err != nil {
		return DatabaseConfig{Path: ""}
	}

	return DatabaseConfig{
		Path: filepath.Join(dir, "hop", "hop.db"),
	}
}
