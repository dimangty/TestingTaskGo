package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"module-name/config"
)

func TestNewRequestAddsMasterKey(t *testing.T) {
	client := New(config.Config{Key: "test-key"})

	request, err := client.NewRequest(context.Background(), http.MethodGet, "/b/bin-id", nil)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	if got := request.Header.Get(masterKeyHeader); got != "test-key" {
		t.Fatalf("%s header = %q, want %q", masterKeyHeader, got, "test-key")
	}
	if got, want := request.URL.String(), baseURL+"/b/bin-id"; got != want {
		t.Fatalf("request URL = %q, want %q", got, want)
	}
}

func TestCreateSendsDocumentAndReturnsMetadata(t *testing.T) {
	createdAt := time.Date(2026, time.September, 2, 12, 30, 0, 0, time.UTC)
	client := testClient(httpClientFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodPost {
			t.Errorf("method = %q, want %q", request.Method, http.MethodPost)
		}
		if request.URL.Path != "/b" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/b")
		}
		if got := request.Header.Get(masterKeyHeader); got != "test-key" {
			t.Errorf("%s = %q, want %q", masterKeyHeader, got, "test-key")
		}
		if got := request.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		if got := request.Header.Get(binNameHeader); got != "my-bin" {
			t.Errorf("%s = %q, want my-bin", binNameHeader, got)
		}
		if got := request.Header.Get(binPrivateHeader); got != "true" {
			t.Errorf("%s = %q, want true", binPrivateHeader, got)
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
		}
		if got, want := string(body), `{"hello":"world"}`; got != want {
			t.Errorf("body = %q, want %q", got, want)
		}

		return jsonResponse(http.StatusOK, `{"record":{"hello":"world"},"metadata":{"id":"bin-id","private":true,"createdAt":"`+createdAt.Format(time.RFC3339)+`"}}`), nil
	}))

	created, err := client.Create(context.Background(), []byte(`{"hello":"world"}`), "my-bin")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID != "bin-id" || created.Name != "my-bin" || !created.Private || !created.CreatedAt.Equal(createdAt) {
		t.Fatalf("Create() = %+v, want response metadata and name", created)
	}
}

func TestGetUpdateAndDeleteUseBinRoute(t *testing.T) {
	wantedMethods := []string{http.MethodGet, http.MethodPut, http.MethodDelete}
	requestIndex := 0
	client := testClient(httpClientFunc(func(request *http.Request) (*http.Response, error) {
		if requestIndex >= len(wantedMethods) {
			t.Fatalf("unexpected extra request: %s %s", request.Method, request.URL.Path)
		}
		if got, want := request.Method, wantedMethods[requestIndex]; got != want {
			t.Errorf("request %d method = %q, want %q", requestIndex, got, want)
		}
		if got, want := request.URL.Path, "/b/bin-id"; got != want {
			t.Errorf("request %d path = %q, want %q", requestIndex, got, want)
		}
		if request.Method == http.MethodPut {
			if got := request.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
		}
		requestIndex++
		return jsonResponse(http.StatusOK, `{"record":{"value":2},"metadata":{"id":"bin-id"}}`), nil
	}))

	record, err := client.Get(context.Background(), "bin-id")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got, want := string(record), `{"value":2}`; got != want {
		t.Fatalf("Get() = %s, want record %s", got, want)
	}
	if _, err := client.Update(context.Background(), "bin-id", []byte(`{"value":2}`)); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if err := client.Delete(context.Background(), "bin-id"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if requestIndex != len(wantedMethods) {
		t.Fatalf("request count = %d, want %d", requestIndex, len(wantedMethods))
	}
}

func TestGetRejectsResponseWithoutRecord(t *testing.T) {
	client := testClient(httpClientFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{"metadata":{"id":"bin-id"}}`), nil
	}))

	_, err := client.Get(context.Background(), "bin-id")
	if err == nil || !strings.Contains(err.Error(), "does not contain a record") {
		t.Fatalf("Get() error = %v, want missing record error", err)
	}
}

func TestGetRejectsMalformedResponse(t *testing.T) {
	client := testClient(httpClientFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{"record":`), nil
	}))

	_, err := client.Get(context.Background(), "bin-id")
	if err == nil || !strings.Contains(err.Error(), "decode response") {
		t.Fatalf("Get() error = %v, want decode response error", err)
	}
}

func TestAPIErrorIncludesStatusAndMessage(t *testing.T) {
	client := testClient(httpClientFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusUnauthorized, `{"message":"invalid key"}`), nil
	}))

	_, err := client.Get(context.Background(), "bin-id")
	var apiError *APIError
	if !errors.As(err, &apiError) {
		t.Fatalf("Get() error = %v, want *APIError", err)
	}
	if apiError.StatusCode != http.StatusUnauthorized || apiError.Message != "invalid key" {
		t.Fatalf("APIError = %+v", apiError)
	}
}

func TestCreateRejectsInvalidJSONBeforeRequest(t *testing.T) {
	client := testClient(httpClientFunc(func(*http.Request) (*http.Response, error) {
		panic("HTTP request must not be sent")
	}))

	_, err := client.Create(context.Background(), []byte(`{"broken"`), "my-bin")
	if !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("Create() error = %v, want %v", err, ErrInvalidJSON)
	}
}

func testClient(httpClient HTTPClient) *Client {
	return &Client{
		key:        "test-key",
		baseURL:    "https://example.test",
		httpClient: httpClient,
	}
}

func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

type httpClientFunc func(*http.Request) (*http.Response, error)

func (function httpClientFunc) Do(request *http.Request) (*http.Response, error) {
	return function(request)
}
