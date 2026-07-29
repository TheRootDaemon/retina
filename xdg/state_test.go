package xdg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserStateDir(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name    string
		goos    string
		env     map[string]string
		want    string
		wantErr bool
	}{
		{
			name: "linux XDG_STATE_HOME set",
			goos: "linux",
			env:  map[string]string{"XDG_STATE_HOME": "/custom/state"},
			want: "/custom/state",
		},
		{
			name: "linux XDG_STATE_HOME empty falls back to home/.local/state",
			goos: "linux",
			env:  map[string]string{"XDG_STATE_HOME": ""},
			want: filepath.Join(home, ".local", "state"),
		},
		{
			name: "linux XDG_STATE_HOME with trailing slash",
			goos: "linux",
			env:  map[string]string{"XDG_STATE_HOME": "/tmp/state/"},
			want: "/tmp/state/",
		},
		{
			name: "linux XDG_STATE_HOME with spaces",
			goos: "linux",
			env:  map[string]string{"XDG_STATE_HOME": "/path with spaces/state"},
			want: "/path with spaces/state",
		},
		{
			name:    "linux no home dir",
			goos:    "linux",
			env:     map[string]string{"XDG_STATE_HOME": "", "HOME": ""},
			wantErr: true,
		},
		{
			name: "darwin default",
			goos: "darwin",
			want: filepath.Join(home, "Library", "Application Support"),
		},
		{
			name:    "darwin no home dir",
			goos:    "darwin",
			env:     map[string]string{"HOME": ""},
			wantErr: true,
		},
		{
			name: "windows LocalAppData set",
			goos: "windows",
			env:  map[string]string{"LocalAppData": `C:\Users\test\AppData\Local`},
			want: `C:\Users\test\AppData\Local`,
		},
		{
			name:    "windows LocalAppData empty",
			goos:    "windows",
			env:     map[string]string{"LocalAppData": ""},
			wantErr: true,
		},
		{
			name: "freebsd XDG_STATE_HOME set",
			goos: "freebsd",
			env:  map[string]string{"XDG_STATE_HOME": "/bsd/state"},
			want: "/bsd/state",
		},
		{
			name: "openbsd XDG_STATE_HOME set",
			goos: "openbsd",
			env:  map[string]string{"XDG_STATE_HOME": "/bsd/state"},
			want: "/bsd/state",
		},
		{
			name: "returns absolute path",
			goos: "linux",
			env:  map[string]string{"XDG_STATE_HOME": "/tmp"},
			want: "/tmp",
		},
		{
			name: "consistent across calls",
			goos: "linux",
			env:  map[string]string{"XDG_STATE_HOME": ""},
			want: filepath.Join(home, ".local", "state"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got, err := userStateDir(tt.goos)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, got)
				assert.Equal(t, tt.want, got)
				if tt.goos != "windows" {
					assert.True(t, filepath.IsAbs(got))
				}
			}
		})
	}
}
