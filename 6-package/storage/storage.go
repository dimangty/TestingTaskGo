// Package storage persists bins in local JSON files.
package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"module-name/bins"
	fileutil "module-name/file"
)

// ErrNotJSON is returned when a storage path does not have a JSON extension.
var ErrNotJSON = errors.New("storage file must have a .json extension")

// SaveBin adds a bin to a local JSON file. If a bin with the same ID already
// exists, it is replaced instead of being duplicated.
func SaveBin(path string, bin bins.Bin) error {
	if !fileutil.IsJSON(path) {
		return fmt.Errorf("save bin: %w", ErrNotJSON)
	}

	storedBins, err := ReadBins(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("save bin: %w", err)
	}

	replaced := false
	for index := range storedBins {
		if storedBins[index].ID == bin.ID {
			storedBins[index] = bin
			replaced = true
			break
		}
	}

	if !replaced {
		storedBins = append(storedBins, bin)
	}

	data, err := json.MarshalIndent(storedBins, "", "  ")
	if err != nil {
		return fmt.Errorf("save bin: encode JSON: %w", err)
	}

	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("save bin: write file %q: %w", path, err)
	}

	return nil
}

// ReadBins reads a list of bins from a local JSON file.
func ReadBins(path string) ([]bins.Bin, error) {
	if !fileutil.IsJSON(path) {
		return nil, fmt.Errorf("read bins: %w", ErrNotJSON)
	}

	data, err := fileutil.Read(path)
	if err != nil {
		return nil, fmt.Errorf("read bins: %w", err)
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return []bins.Bin{}, nil
	}

	var storedBins []bins.Bin
	if err := json.Unmarshal(data, &storedBins); err != nil {
		return nil, fmt.Errorf("read bins: decode JSON from %q: %w", path, err)
	}

	if storedBins == nil {
		return []bins.Bin{}, nil
	}

	return storedBins, nil
}
