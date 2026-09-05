// Package storage persists bins in local JSON files.
package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"

	"module-name/bins"
)

var (
	// ErrNotJSON is returned when a storage path does not have a JSON extension.
	ErrNotJSON = errors.New("storage file must have a .json extension")
	// ErrFileRequired is returned when storage is created without a file dependency.
	ErrFileRequired = errors.New("file dependency is required")
)

// File describes the file operations required by JSONStorage.
// The concrete implementation is injected by the application entry point.
type File interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
	IsJSON(path string) bool
}

// JSONStorage persists bins as a JSON array.
type JSONStorage struct {
	path string
	file File
}

// NewJSONStorage creates a JSON-backed bin storage.
func NewJSONStorage(path string, file File) *JSONStorage {
	return &JSONStorage{
		path: path,
		file: file,
	}
}

// SaveBin adds a bin to storage. If a bin with the same ID already exists, it
// is replaced instead of being duplicated.
func (storage *JSONStorage) SaveBin(bin bins.Bin) error {
	if storage == nil || storage.file == nil {
		return fmt.Errorf("save bin: %w", ErrFileRequired)
	}
	if !storage.file.IsJSON(storage.path) {
		return fmt.Errorf("save bin: %w", ErrNotJSON)
	}

	storedBins, err := storage.ReadBins()
	if err != nil {
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

	if err := storage.writeBins(storedBins); err != nil {
		return fmt.Errorf("save bin: %w", err)
	}

	return nil
}

// DeleteBin removes a bin from storage. A missing ID is treated as success.
func (storage *JSONStorage) DeleteBin(id string) error {
	if storage == nil || storage.file == nil {
		return fmt.Errorf("delete bin: %w", ErrFileRequired)
	}
	if !storage.file.IsJSON(storage.path) {
		return fmt.Errorf("delete bin: %w", ErrNotJSON)
	}

	storedBins, err := storage.ReadBins()
	if err != nil {
		return fmt.Errorf("delete bin: %w", err)
	}

	for index := range storedBins {
		if storedBins[index].ID == id {
			storedBins = append(storedBins[:index], storedBins[index+1:]...)
			if err := storage.writeBins(storedBins); err != nil {
				return fmt.Errorf("delete bin: %w", err)
			}
			return nil
		}
	}

	return nil
}

// ReadBins reads all bins from storage. A missing file is treated as an empty
// storage so the application can start without a separate initialization step.
func (storage *JSONStorage) ReadBins() ([]bins.Bin, error) {
	if storage == nil || storage.file == nil {
		return nil, fmt.Errorf("read bins: %w", ErrFileRequired)
	}
	if !storage.file.IsJSON(storage.path) {
		return nil, fmt.Errorf("read bins: %w", ErrNotJSON)
	}

	data, err := storage.file.Read(storage.path)
	if errors.Is(err, fs.ErrNotExist) {
		return []bins.Bin{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read bins: %w", err)
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return []bins.Bin{}, nil
	}

	var storedBins []bins.Bin
	if err := json.Unmarshal(data, &storedBins); err != nil {
		return nil, fmt.Errorf("read bins: decode JSON from %q: %w", storage.path, err)
	}

	if storedBins == nil {
		return []bins.Bin{}, nil
	}

	return storedBins, nil
}

func (storage *JSONStorage) writeBins(storedBins []bins.Bin) error {
	data, err := json.MarshalIndent(storedBins, "", "  ")
	if err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	data = append(data, '\n')
	if err := storage.file.Write(storage.path, data); err != nil {
		return err
	}
	return nil
}
