package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestC_Defaults(t *testing.T) {
	ResetForTesting()
	defer ResetForTesting()

	cfg := C()
	require.NotNil(t, cfg)
	assert.Equal(t, Default(), *cfg)
}

func TestC_ReturnsSamePointer(t *testing.T) {
	ResetForTesting()
	defer ResetForTesting()

	d := Default()
	currentConfig.Store(&d)

	cfg1 := C()
	cfg2 := C()
	assert.Same(t, cfg1, cfg2)
}

func TestInitialize_and_C(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte(`[logger]
enabled = true
level = "debug"
`), 0o600)
	require.NoError(t, err)

	ResetForTesting()
	defer ResetForTesting()

	t.Setenv("HOP_CONFIG", path)
	err = Initialize()
	require.NoError(t, err)

	cfg := C()
	assert.Equal(t, "debug", cfg.Logger.Level)
	assert.True(t, cfg.Logger.Enabled)
}

func TestInitialize_withDatabaseOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte(`[database]
path = ":memory:"
`), 0o600)
	require.NoError(t, err)

	ResetForTesting()
	defer ResetForTesting()

	t.Setenv("HOP_CONFIG", path)
	err = Initialize()
	require.NoError(t, err)

	assert.Equal(t, ":memory:", Database().Path)
}

func TestInitialize_MissingFileDefaults(t *testing.T) {
	ResetForTesting()
	defer ResetForTesting()

	t.Setenv("HOP_CONFIG", "/nonexistent/path/config.toml")
	err := Initialize()
	require.NoError(t, err)

	cfg := C()
	require.NotNil(t, cfg)
	assert.Equal(t, Default(), *cfg)
}

func TestInitialize_MalformedFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte("invalid toml content {{{"), 0o600)
	require.NoError(t, err)

	ResetForTesting()
	defer ResetForTesting()

	t.Setenv("HOP_CONFIG", path)
	err = Initialize()
	require.Error(t, err)
}

func TestInitialize_Idempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte("[logger]\nlevel = \"warn\"\n"), 0o600)
	require.NoError(t, err)

	ResetForTesting()
	defer ResetForTesting()

	t.Setenv("HOP_CONFIG", path)

	err = Initialize()
	require.NoError(t, err)
	err = Initialize()
	require.NoError(t, err)

	cfg := C()
	assert.Equal(t, "warn", cfg.Logger.Level)
}

func TestInitialize_SaltAndVault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte("salt = \"my-salt\"\nvault = \"/data/vault\"\n"), 0o600)
	require.NoError(t, err)

	ResetForTesting()
	defer ResetForTesting()

	t.Setenv("HOP_CONFIG", path)
	err = Initialize()
	require.NoError(t, err)

	cfg := C()
	assert.Equal(t, "my-salt", cfg.Salt)
	assert.Equal(t, "/data/vault", cfg.Vault)
}

func TestInitialize_ReplacesPreviousConfig(t *testing.T) {
	dir := t.TempDir()
	path1 := filepath.Join(dir, "config1.toml")
	path2 := filepath.Join(dir, "config2.toml")
	require.NoError(t, os.WriteFile(path1, []byte("[logger]\nlevel = \"debug\"\n"), 0o600))
	require.NoError(t, os.WriteFile(path2, []byte("[logger]\nlevel = \"error\"\n"), 0o600))

	ResetForTesting()
	defer ResetForTesting()

	t.Setenv("HOP_CONFIG", path1)
	require.NoError(t, Initialize())
	assert.Equal(t, "debug", C().Logger.Level)

	t.Setenv("HOP_CONFIG", path2)
	require.NoError(t, Initialize())
	assert.Equal(t, "error", C().Logger.Level)
}

func TestLogger_Accessor(t *testing.T) {
	ResetForTesting()
	defer ResetForTesting()

	lc := Logger()
	assert.Equal(t, DefaultLoggerConfig(), lc)
}

func TestLogger_AccessorAfterInitialize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	require.NoError(t, os.WriteFile(path, []byte("[logger]\nenabled = true\nlevel = \"trace\"\n"), 0o600))

	ResetForTesting()
	defer ResetForTesting()

	t.Setenv("HOP_CONFIG", path)
	require.NoError(t, Initialize())

	lc := Logger()
	assert.True(t, lc.Enabled)
	assert.Equal(t, "trace", lc.Level)
}

func TestDatabase_Accessor(t *testing.T) {
	ResetForTesting()
	defer ResetForTesting()

	dc := Database()
	assert.Equal(t, DefaultDatabaseConfig(), dc)
}

func TestDatabase_AccessorAfterInitialize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	require.NoError(t, os.WriteFile(path, []byte("[database]\npath = \"/custom/db.sqlite\"\n"), 0o600))

	ResetForTesting()
	defer ResetForTesting()

	t.Setenv("HOP_CONFIG", path)
	require.NoError(t, Initialize())

	dc := Database()
	assert.Equal(t, "/custom/db.sqlite", dc.Path)
}

func TestResetForTesting(t *testing.T) {
	ResetForTesting()
	defer ResetForTesting()

	d := Default()
	currentConfig.Store(&d)
	assert.NotNil(t, currentConfig.Load())

	ResetForTesting()
	assert.Nil(t, currentConfig.Load())
}

func TestResetForTesting_ThenCReturnsDefaults(t *testing.T) {
	ResetForTesting()
	defer ResetForTesting()

	d := Default()
	currentConfig.Store(&d)

	ResetForTesting()
	cfg := C()
	assert.Equal(t, Default(), *cfg)
}
