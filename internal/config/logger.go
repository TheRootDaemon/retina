package config

import (
	"os"
	"path/filepath"

	"github.com/therootdaemon/hop/xdg"
)

// LoggerConfig holds logger-specific configuration.
// It maps to the [logger] section in the TOML config file.
//
// Example TOML:
//
//	[logger]
//	  Level = "info"
//	  File  = "hop/hop.log"
type LoggerConfig struct {
	Enabled bool
	// Level is the minimum log level ("trace", "debug", "info", "warn", "error").
	Level string `toml:"level"`

	// File is the log file path relative to UserStateDir.
	// Defaults to "hop/hop.log".
	File string `toml:"file"`
}

// DefaultLoggerConfig returns the default logger settings.
func DefaultLoggerConfig() LoggerConfig {
	file, _ := DefaultLogFile()

	return LoggerConfig{
		Level: "info",
		File:  file,
	}
}

// DefaultLogFile returns the default absolute path to the log file.
//
// The path is located beneath the user's state directory:
//
//   - Linux/BSD:  $XDG_STATE_HOME/hop/hop.log or $HOME/.local/state/hop/hop.log
//   - macOS:      $HOME/Library/Application Support/hop/hop.log
//   - Windows:    %LocalAppData%\hop\hop.log
//
// If the state directory cannot be determined, it falls back to
// FallbackLogFile.
func DefaultLogFile() (string, error) {
	stateDir, err := xdg.UserStateDir()
	if err != nil {
		return FallbackLogFile()
	}
	return filepath.Join(stateDir, "hop", "hop.log"), nil
}

// FallbackLogFile returns the fallback absolute path to the log file.
//
// The fallback location is "$HOME/.hop/hop.log". It is used when the
// platform-specific state directory cannot be determined.
func FallbackLogFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".hop", "hop.log"), nil
}
