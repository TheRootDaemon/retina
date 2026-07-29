package xdg

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// UserStateDir returns the directory used to store user-specific state data,
// such as logs, history, caches
// that must persist across restarts, and undo information.
//
// The returned directory depends on the current operating system:
//
//   - On Unix-like systems, it returns the value of the XDG_STATE_HOME
//     environment variable if set; otherwise it falls back to
//     $HOME/.local/state, as defined by the XDG Base Directory Specification.
//
//   - On macOS, it returns $HOME/Library/Application Support.
//
//   - On Windows, it returns the value of the LocalAppData environment
//     variable.
//
// An error is returned if the user's home directory cannot be determined or if
// the required platform-specific environment variable is unavailable.
func UserStateDir() (string, error) {
	return userStateDir(runtime.GOOS)
}

// userStateDir returns the platform-specific user state directory
// for the provided operating system name.
func userStateDir(goos string) (string, error) {
	switch goos {
	case "windows":
		if dir := os.Getenv("LocalAppData"); dir != "" {
			return dir, nil
		}
		return "", errors.New("%LocalAppData% is not defined")

	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support"), nil

	default: // Linux, BSD, etc.
		if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
			return dir, nil
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".local", "state"), nil
	}
}
