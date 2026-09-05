package bins

import (
	"errors"
	"testing"
)

func TestBinListAddReplacesExistingBin(t *testing.T) {
	store := &storageStub{bins: []Bin{{ID: "one", Name: "old"}}}
	list, err := NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}

	if err := list.Add(Bin{ID: "one", Name: "new"}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if got := list.Bins(); len(got) != 1 || got[0].Name != "new" {
		t.Fatalf("Bins() = %+v, want one replacement", got)
	}
	if len(store.saved) != 1 || store.saved[0].Name != "new" {
		t.Fatalf("saved bins = %+v, want replacement", store.saved)
	}
}

func TestBinListDeletePersistsAndRemovesBin(t *testing.T) {
	store := &storageStub{bins: []Bin{{ID: "one"}, {ID: "two"}}}
	list, err := NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}

	if err := list.Delete("one"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if store.deletedID != "one" {
		t.Fatalf("deleted ID = %q, want one", store.deletedID)
	}
	if got := list.Bins(); len(got) != 1 || got[0].ID != "two" {
		t.Fatalf("Bins() = %+v, want only bin two", got)
	}
}

func TestBinListDoesNotChangeMemoryWhenStorageFails(t *testing.T) {
	storageError := errors.New("storage unavailable")
	store := &storageStub{
		bins:      []Bin{{ID: "one"}},
		deleteErr: storageError,
	}
	list, err := NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}

	err = list.Delete("one")
	if !errors.Is(err, storageError) {
		t.Fatalf("Delete() error = %v, want %v", err, storageError)
	}
	if got := list.Bins(); len(got) != 1 || got[0].ID != "one" {
		t.Fatalf("Bins() = %+v, want original list", got)
	}
}

type storageStub struct {
	bins      []Bin
	saved     []Bin
	deletedID string
	deleteErr error
}

func (storage *storageStub) ReadBins() ([]Bin, error) {
	return append([]Bin(nil), storage.bins...), nil
}

func (storage *storageStub) SaveBin(bin Bin) error {
	storage.saved = append(storage.saved, bin)
	return nil
}

func (storage *storageStub) DeleteBin(id string) error {
	storage.deletedID = id
	return storage.deleteErr
}
