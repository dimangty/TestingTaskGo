// Package api contains code for communicating with the JSON.BIN API.
package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"module-name/config"
)

const (
	baseURL         = "https://api.jsonbin.io/v3"
	masterKeyHeader = "X-Master-Key"
)

// Client creates authenticated requests to JSON.BIN.
type Client struct {
	key string
}

// New creates a JSON.BIN API client from the application configuration.
func New(appConfig config.Config) *Client {
	return &Client{key: appConfig.Key}
}

// NewRequest creates a request and adds the JSON.BIN master key header.
func (client *Client) NewRequest(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
) (*http.Request, error) {
	if client == nil || strings.TrimSpace(client.key) == "" {
		return nil, config.ErrKeyRequired
	}

	requestURL := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
	request, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("create JSON.BIN request: %w", err)
	}

	request.Header.Set(masterKeyHeader, client.key)
	return request, nil
}
