// Package file contains helpers for working with local files.
package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Local reads and writes files on the local filesystem.
// It is intentionally stateless, so a single value can be shared by services.
type Local struct{}

// NewLocal creates a local filesystem implementation.
func NewLocal() Local {
	return Local{}
}

// Read returns the contents of a local file.
func (Local) Read(path string) ([]byte, error) {
	return Read(path)
}

// Write replaces a local file with data.
func (Local) Write(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write file %q: %w", path, err)
	}

	return nil
}

// IsJSON reports whether path has a .json extension.
func (Local) IsJSON(path string) bool {
	return IsJSON(path)
}

// Read reads and returns the contents of a local file.
func Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", path, err)
	}

	return data, nil
}

// IsJSON reports whether path has a .json extension.
// The check is case-insensitive.
func IsJSON(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".json")
}
