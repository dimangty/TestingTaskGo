package storage

import (
	"path/filepath"
	"reflect"
	"testing"

	"module-name/bins"
	localfile "module-name/file"
)

func TestJSONStorageSavesBin(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	storagePath := filepath.Join(t.TempDir(), "bins.json")
	store := NewJSONStorage(storagePath, localfile.NewLocal())
	expected := []bins.Bin{{ID: "one", Name: "first"}}

	// Act - выполняем функцию
	err := store.SaveBin(expected[0])

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("SaveBin() error = %v", err)
	}
	got, err := store.ReadBins()
	if err != nil {
		t.Fatalf("ReadBins() error = %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("ReadBins() = %+v, want %+v", got, expected)
	}
}

func TestJSONStorageReplacesBin(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	storagePath := filepath.Join(t.TempDir(), "bins.json")
	store := NewJSONStorage(storagePath, localfile.NewLocal())
	if err := store.SaveBin(bins.Bin{ID: "one", Name: "first"}); err != nil {
		t.Fatalf("arrange storage: %v", err)
	}
	expected := []bins.Bin{{ID: "one", Name: "renamed"}}

	// Act - выполняем функцию
	err := store.SaveBin(expected[0])

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("SaveBin() error = %v", err)
	}
	got, err := store.ReadBins()
	if err != nil {
		t.Fatalf("ReadBins() error = %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("ReadBins() = %+v, want %+v", got, expected)
	}
}

func TestJSONStorageDeletesBin(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	storagePath := filepath.Join(t.TempDir(), "bins.json")
	store := NewJSONStorage(storagePath, localfile.NewLocal())
	if err := store.SaveBin(bins.Bin{ID: "one", Name: "first"}); err != nil {
		t.Fatalf("arrange first bin: %v", err)
	}
	if err := store.SaveBin(bins.Bin{ID: "two", Name: "second"}); err != nil {
		t.Fatalf("arrange second bin: %v", err)
	}
	expected := []bins.Bin{{ID: "one", Name: "first"}}

	// Act - выполняем функцию
	err := store.DeleteBin("two")

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("DeleteBin() error = %v", err)
	}
	got, err := store.ReadBins()
	if err != nil {
		t.Fatalf("ReadBins() error = %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("ReadBins() = %+v, want %+v", got, expected)
	}
}

func TestJSONStorageDeleteMissingBinIsIdempotent(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	storagePath := filepath.Join(t.TempDir(), "bins.json")
	store := NewJSONStorage(storagePath, localfile.NewLocal())
	expected := []bins.Bin{}

	// Act - выполняем функцию
	err := store.DeleteBin("missing")

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("DeleteBin(missing) error = %v", err)
	}
	got, err := store.ReadBins()
	if err != nil {
		t.Fatalf("ReadBins() error = %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("ReadBins() = %+v, want %+v", got, expected)
	}
}
