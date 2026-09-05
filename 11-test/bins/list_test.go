package bins

import (
	"errors"
	"reflect"
	"testing"
)

func TestBinListAddReplacesExistingBin(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	store := &storageStub{bins: []Bin{{ID: "one", Name: "old"}}}
	list, err := NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Bin{{ID: "one", Name: "new"}}

	// Act - выполняем функцию
	err = list.Add(expected[0])

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if got := list.Bins(); !reflect.DeepEqual(got, expected) {
		t.Fatalf("Bins() = %+v, want %+v", got, expected)
	}
	if !reflect.DeepEqual(store.saved, expected) {
		t.Fatalf("saved bins = %+v, want %+v", store.saved, expected)
	}
}

func TestBinListDeletePersistsAndRemovesBin(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	store := &storageStub{bins: []Bin{{ID: "one"}, {ID: "two"}}}
	list, err := NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Bin{{ID: "two"}}

	// Act - выполняем функцию
	err = list.Delete("one")

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if store.deletedID != "one" {
		t.Fatalf("deleted ID = %q, want one", store.deletedID)
	}
	if got := list.Bins(); !reflect.DeepEqual(got, expected) {
		t.Fatalf("Bins() = %+v, want %+v", got, expected)
	}
}

func TestBinListDoesNotChangeMemoryWhenStorageFails(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	storageError := errors.New("storage unavailable")
	expected := []Bin{{ID: "one"}}
	store := &storageStub{
		bins:      expected,
		deleteErr: storageError,
	}
	list, err := NewBinList(store)
	if err != nil {
		t.Fatal(err)
	}

	// Act - выполняем функцию
	err = list.Delete("one")

	// Assert - проверка результата с expected
	if !errors.Is(err, storageError) {
		t.Fatalf("Delete() error = %v, want %v", err, storageError)
	}
	if got := list.Bins(); !reflect.DeepEqual(got, expected) {
		t.Fatalf("Bins() = %+v, want %+v", got, expected)
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
