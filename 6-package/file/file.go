// Package file contains helpers for working with local files.
package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
