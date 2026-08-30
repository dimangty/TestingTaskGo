package bins

import (
	"errors"
	"fmt"
)

// ErrStorageRequired is returned when a bin list is created without storage.
var ErrStorageRequired = errors.New("bin storage is required")

// Storage describes the persistence operations required by BinList.
// Implementations are provided by the application entry point.
type Storage interface {
	ReadBins() ([]Bin, error)
	SaveBin(bin Bin) error
}

// BinList contains a list of bins.
type BinList struct {
	bins    []Bin
	storage Storage
}

// NewBinList loads a bin list from the supplied storage.
func NewBinList(storage Storage) (*BinList, error) {
	if storage == nil {
		return nil, ErrStorageRequired
	}

	storedBins, err := storage.ReadBins()
	if err != nil {
		return nil, fmt.Errorf("create bin list: %w", err)
	}

	return &BinList{
		bins:    append([]Bin(nil), storedBins...),
		storage: storage,
	}, nil
}

// Add persists a bin and adds it to the list. A bin with an existing ID is
// replaced instead of being duplicated.
func (list *BinList) Add(bin Bin) error {
	if err := list.storage.SaveBin(bin); err != nil {
		return fmt.Errorf("add bin: %w", err)
	}

	for index := range list.bins {
		if list.bins[index].ID == bin.ID {
			list.bins[index] = bin
			return nil
		}
	}

	list.bins = append(list.bins, bin)
	return nil
}

// Bins returns a copy of the bins in the list.
func (list *BinList) Bins() []Bin {
	return append([]Bin(nil), list.bins...)
}
