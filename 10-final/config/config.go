// Package config loads application settings from environment variables.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const keyEnv = "KEY"

// ErrKeyRequired is returned when the JSON.BIN API key is not configured.
var ErrKeyRequired = errors.New("JSON.BIN API key is required")

// Config contains settings required by the application.
type Config struct {
	Key string
}

// Load reads settings from .env and the process environment.
// Values already present in the process environment take precedence.
func Load() (Config, error) {
	return LoadFrom(".env")
}

// LoadFrom reads settings using the supplied dotenv file. A missing file is
// allowed so the application can also be configured only through its process
// environment.
func LoadFrom(path string) (Config, error) {
	if err := godotenv.Load(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Config{}, fmt.Errorf("load config from %q: %w", path, err)
	}

	key := strings.TrimSpace(os.Getenv(keyEnv))
	if key == "" {
		return Config{}, ErrKeyRequired
	}

	return Config{Key: key}, nil
}
