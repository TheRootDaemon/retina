package salt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	salt, err := Generate()
	require.NoError(t, err)
	assert.Len(t, salt, Size)
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantErr error
	}{
		{
			name: "valid salt",
			setup: func(t *testing.T) string {
				p := filepath.Join(t.TempDir(), "salt.bin")
				require.NoError(t, Save(p, make([]byte, Size)))
				return p
			},
		},
		{
			name: "wrong size",
			setup: func(t *testing.T) string {
				p := filepath.Join(t.TempDir(), "salt.bin")
				require.NoError(t, os.WriteFile(p, make([]byte, 4), 0o600))
				return p
			},
			wantErr: ErrInvalidSalt,
		},
		{
			name: "nonexistent path",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nope.bin")
			},
			wantErr: os.ErrNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)
			_, err := Load(path)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		name    string
		pathFn  func(t *testing.T) string
		size    int
		wantErr error
	}{
		{
			name: "valid salt",
			pathFn: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "salt.bin")
			},
			size: Size,
		},
		{
			name: "wrong size",
			pathFn: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "salt.bin")
			},
			size:    4,
			wantErr: ErrInvalidSalt,
		},
		{
			name: "parent dir created",
			pathFn: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "a", "b", "salt.bin")
			},
			size: Size,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.pathFn(t)
			salt := make([]byte, tt.size)
			err := Save(path, salt)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				got, err := os.ReadFile(path)
				require.NoError(t, err)
				assert.Equal(t, salt, got)
			}
		})
	}
}

func TestLoadOrCreate(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T) string
		verify func(t *testing.T, path string, salt []byte)
	}{
		{
			name: "file exists",
			setup: func(t *testing.T) string {
				p := filepath.Join(t.TempDir(), "salt.bin")
				require.NoError(t, Save(p, make([]byte, Size)))
				return p
			},
			verify: func(t *testing.T, path string, salt []byte) {
				assert.Len(t, salt, Size)
				assert.Equal(t, make([]byte, Size), salt)
			},
		},
		{
			name: "file does not exist",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "salt.bin")
			},
			verify: func(t *testing.T, path string, salt []byte) {
				assert.Len(t, salt, Size)
				_, err := os.Stat(path)
				assert.NoError(t, err, "salt file should have been created")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)
			salt, err := LoadOrCreate(path)
			require.NoError(t, err)
			tt.verify(t, path, salt)
		})
	}
}

func TestLoadOrCreate_NonNotExistError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "salt.bin")
	require.NoError(t, os.WriteFile(path, []byte("too short"), 0o600))

	_, err := LoadOrCreate(path)
	assert.ErrorIs(t, err, ErrInvalidSalt)
}
