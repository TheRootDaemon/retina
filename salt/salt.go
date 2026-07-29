package salt

import (
	"crypto/rand"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Size is the required length, in bytes, of a valid salt.
const Size = 32

// ErrInvalidSalt indicates that a salt is malformed or has an invalid length.
var ErrInvalidSalt = errors.New("invalid salt")

// Generate returns a cryptographically random salt of Size bytes.
func Generate() ([]byte, error) {
	salt := make([]byte, Size)

	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	return salt, nil
}

// Load reads a salt from path
// and validates its length is exactly Size.
//
// The file is opened using os.OpenRoot for path-traversal safety.
// If the salt length does not match, ErrInvalidSalt is returned.
func Load(path string) ([]byte, error) {
	r, file, err := openRoot(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Close()
	}()

	f, err := r.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	salt, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if len(salt) != Size {
		return nil, ErrInvalidSalt
	}

	return salt, nil
}

// Save writes salt to path.
//
// The parent directory is created if necessary.
// The salt must be exactly Size bytes.
func Save(path string, salt []byte) error {
	if len(salt) != Size {
		return ErrInvalidSalt
	}

	if err := os.MkdirAll(
		filepath.Dir(path),
		0o750,
	); err != nil {
		return err
	}

	r, file, err := openRoot(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()

	f, err := r.OpenFile(
		file,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o600,
	)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err = f.Write(salt); err != nil {
		return err
	}

	return f.Sync()
}

// LoadOrCreate loads an existing salt from path.
//
// If the file does not exist, it generates a new salt,
// saves it to path, and returns it.
func LoadOrCreate(path string) ([]byte, error) {
	salt, err := Load(path)
	if err == nil {
		return salt, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	salt, err = Generate()
	if err != nil {
		return nil, err
	}

	if err := Save(path, salt); err != nil {
		return nil, err
	}

	return salt, nil
}

func openRoot(path string) (*os.Root, string, error) {
	root := filepath.Dir(path)
	file := filepath.Base(path)

	r, err := os.OpenRoot(root)
	if err != nil {
		return nil, "", err
	}

	return r, file, nil
}
