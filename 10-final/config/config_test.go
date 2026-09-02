package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromReadsKeyFromDotenv(t *testing.T) {
	previousKey, hadKey := os.LookupEnv(keyEnv)
	if err := os.Unsetenv(keyEnv); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if hadKey {
			_ = os.Setenv(keyEnv, previousKey)
			return
		}
		_ = os.Unsetenv(keyEnv)
	})

	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("KEY=test-key\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}
	if got.Key != "test-key" {
		t.Fatalf("LoadFrom().Key = %q, want %q", got.Key, "test-key")
	}
}

func TestLoadFromReturnsErrorWhenKeyIsEmpty(t *testing.T) {
	t.Setenv(keyEnv, "")

	_, err := LoadFrom(filepath.Join(t.TempDir(), "missing.env"))
	if !errors.Is(err, ErrKeyRequired) {
		t.Fatalf("LoadFrom() error = %v, want %v", err, ErrKeyRequired)
	}
}
