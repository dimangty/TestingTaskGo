package storage

import (
	"path/filepath"
	"testing"

	"module-name/bins"
	localfile "module-name/file"
)

func TestJSONStorageSavesReplacesAndDeletesBins(t *testing.T) {
	storagePath := filepath.Join(t.TempDir(), "bins.json")
	store := NewJSONStorage(storagePath, localfile.NewLocal())

	if err := store.SaveBin(bins.Bin{ID: "one", Name: "first"}); err != nil {
		t.Fatalf("SaveBin(first) error = %v", err)
	}
	if err := store.SaveBin(bins.Bin{ID: "two", Name: "second"}); err != nil {
		t.Fatalf("SaveBin(second) error = %v", err)
	}
	if err := store.SaveBin(bins.Bin{ID: "one", Name: "renamed"}); err != nil {
		t.Fatalf("SaveBin(replacement) error = %v", err)
	}
	if err := store.DeleteBin("two"); err != nil {
		t.Fatalf("DeleteBin() error = %v", err)
	}

	got, err := store.ReadBins()
	if err != nil {
		t.Fatalf("ReadBins() error = %v", err)
	}
	if len(got) != 1 || got[0].ID != "one" || got[0].Name != "renamed" {
		t.Fatalf("ReadBins() = %+v, want only renamed bin one", got)
	}
}

func TestJSONStorageDeleteMissingBinIsIdempotent(t *testing.T) {
	storagePath := filepath.Join(t.TempDir(), "bins.json")
	store := NewJSONStorage(storagePath, localfile.NewLocal())

	if err := store.DeleteBin("missing"); err != nil {
		t.Fatalf("DeleteBin(missing) error = %v", err)
	}
}
