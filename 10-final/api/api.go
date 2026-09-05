// Package api contains code for communicating with the JSONBin.io API.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"module-name/bins"
	"module-name/config"
)

const (
	baseURL          = "https://api.jsonbin.io/v3"
	masterKeyHeader  = "X-Master-Key"
	binNameHeader    = "X-Bin-Name"
	binPrivateHeader = "X-Bin-Private"
	maxResponseSize  = 2 << 20
)

var (
	// ErrIDRequired is returned when an operation needs a bin ID.
	ErrIDRequired = errors.New("bin ID is required")
	// ErrInvalidJSON is returned before a malformed JSON document is sent.
	ErrInvalidJSON = errors.New("invalid JSON document")
)

// HTTPClient is the part of http.Client used by Client.
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

// Client communicates with JSONBin.io using authenticated requests.
type Client struct {
	key        string
	baseURL    string
	httpClient HTTPClient
}

// APIError describes a non-success response returned by JSONBin.io.
type APIError struct {
	StatusCode int
	Message    string
}

func (err *APIError) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("JSONBin.io returned HTTP %d", err.StatusCode)
	}
	return fmt.Sprintf("JSONBin.io returned HTTP %d: %s", err.StatusCode, err.Message)
}

// New creates a JSONBin.io API client from the application configuration.
func New(appConfig config.Config) *Client {
	return &Client{
		key:     appConfig.Key,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewRequest creates a request and adds the JSONBin.io master key header.
func (client *Client) NewRequest(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
) (*http.Request, error) {
	if client == nil || strings.TrimSpace(client.key) == "" {
		return nil, config.ErrKeyRequired
	}

	requestURL := strings.TrimRight(client.baseURL, "/") + "/" + strings.TrimLeft(path, "/")
	request, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("create JSONBin.io request: %w", err)
	}

	request.Header.Set(masterKeyHeader, client.key)
	return request, nil
}

// Create uploads a JSON document as a private bin and returns its metadata.
func (client *Client) Create(ctx context.Context, document []byte, name string) (bins.Bin, error) {
	if !json.Valid(document) {
		return bins.Bin{}, ErrInvalidJSON
	}

	request, err := client.NewRequest(ctx, http.MethodPost, "/b", bytes.NewReader(document))
	if err != nil {
		return bins.Bin{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(binPrivateHeader, "true")
	if strings.TrimSpace(name) != "" {
		request.Header.Set(binNameHeader, name)
	}

	responseBody, err := client.do(request)
	if err != nil {
		return bins.Bin{}, fmt.Errorf("create bin: %w", err)
	}

	var response struct {
		Metadata struct {
			ID        string    `json:"id"`
			Private   bool      `json:"private"`
			CreatedAt time.Time `json:"createdAt"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return bins.Bin{}, fmt.Errorf("create bin: decode response: %w", err)
	}
	if strings.TrimSpace(response.Metadata.ID) == "" {
		return bins.Bin{}, errors.New("create bin: response does not contain a bin ID")
	}

	return bins.Bin{
		ID:        response.Metadata.ID,
		Private:   response.Metadata.Private,
		CreatedAt: response.Metadata.CreatedAt,
		Name:      name,
	}, nil
}

// Get downloads a bin by ID and returns the stored JSON document.
func (client *Client) Get(ctx context.Context, id string) ([]byte, error) {
	request, err := client.requestForID(ctx, http.MethodGet, id, nil)
	if err != nil {
		return nil, err
	}

	responseBody, err := client.do(request)
	if err != nil {
		return nil, fmt.Errorf("get bin: %w", err)
	}

	var response struct {
		Record json.RawMessage `json:"record"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("get bin: decode response: %w", err)
	}
	if len(response.Record) == 0 {
		return nil, errors.New("get bin: response does not contain a record")
	}

	return response.Record, nil
}

// Update replaces the JSON document stored in a bin.
func (client *Client) Update(ctx context.Context, id string, document []byte) ([]byte, error) {
	if !json.Valid(document) {
		return nil, ErrInvalidJSON
	}

	request, err := client.requestForID(ctx, http.MethodPut, id, bytes.NewReader(document))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	responseBody, err := client.do(request)
	if err != nil {
		return nil, fmt.Errorf("update bin: %w", err)
	}
	return responseBody, nil
}

// Delete removes a bin and all its versions.
func (client *Client) Delete(ctx context.Context, id string) error {
	request, err := client.requestForID(ctx, http.MethodDelete, id, nil)
	if err != nil {
		return err
	}
	if _, err := client.do(request); err != nil {
		return fmt.Errorf("delete bin: %w", err)
	}
	return nil
}

func (client *Client) requestForID(
	ctx context.Context,
	method string,
	id string,
	body io.Reader,
) (*http.Request, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrIDRequired
	}
	return client.NewRequest(ctx, method, "/b/"+url.PathEscape(id), body)
}

func (client *Client) do(request *http.Request) ([]byte, error) {
	if client == nil || client.httpClient == nil {
		return nil, errors.New("HTTP client is required")
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer response.Body.Close()

	limitedBody := io.LimitReader(response.Body, maxResponseSize+1)
	responseBody, err := io.ReadAll(limitedBody)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if len(responseBody) > maxResponseSize {
		return nil, fmt.Errorf("response is larger than %d bytes", maxResponseSize)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		var errorResponse struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(responseBody, &errorResponse)
		return nil, &APIError{
			StatusCode: response.StatusCode,
			Message:    errorResponse.Message,
		}
	}

	return responseBody, nil
}
