package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromReadsKeyFromDotenv(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	expected := "test-key"
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
	if err := os.WriteFile(path, []byte("KEY="+expected+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Act - выполняем функцию
	got, err := LoadFrom(path)

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}
	if got.Key != expected {
		t.Fatalf("LoadFrom().Key = %q, want %q", got.Key, expected)
	}
}

func TestLoadFromReturnsErrorWhenKeyIsEmpty(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	expected := ErrKeyRequired
	t.Setenv(keyEnv, "")
	path := filepath.Join(t.TempDir(), "missing.env")

	// Act - выполняем функцию
	_, err := LoadFrom(path)

	// Assert - проверка результата с expected
	if !errors.Is(err, expected) {
		t.Fatalf("LoadFrom() error = %v, want %v", err, expected)
	}
}
