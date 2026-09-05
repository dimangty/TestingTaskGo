package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"module-name/bins"
	"module-name/config"
)

func TestNewRequestAddsMasterKey(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	expectedKey := "test-key"
	expectedURL := baseURL + "/b/bin-id"
	client := New(config.Config{Key: "test-key"})

	// Act - выполняем функцию
	request, err := client.NewRequest(context.Background(), http.MethodGet, "/b/bin-id", nil)

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	if got := request.Header.Get(masterKeyHeader); got != expectedKey {
		t.Fatalf("%s header = %q, want %q", masterKeyHeader, got, expectedKey)
	}
	if got := request.URL.String(); got != expectedURL {
		t.Fatalf("request URL = %q, want %q", got, expectedURL)
	}
}

func TestCreateSendsDocumentAndReturnsMetadata(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	createdAt := time.Date(2026, time.September, 2, 12, 30, 0, 0, time.UTC)
	expected := bins.Bin{ID: "bin-id", Private: true, CreatedAt: createdAt, Name: "my-bin"}
	var capturedRequest *http.Request
	var capturedBody []byte
	var capturedBodyError error
	client := testClient(httpClientFunc(func(request *http.Request) (*http.Response, error) {
		capturedRequest = request
		capturedBody, capturedBodyError = io.ReadAll(request.Body)
		return jsonResponse(http.StatusOK, `{"record":{"hello":"world"},"metadata":{"id":"bin-id","private":true,"createdAt":"`+createdAt.Format(time.RFC3339)+`"}}`), nil
	}))

	// Act - выполняем функцию
	created, err := client.Create(context.Background(), []byte(`{"hello":"world"}`), "my-bin")

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if capturedRequest == nil {
		t.Fatal("Create() did not send an HTTP request")
	}
	if capturedRequest.Method != http.MethodPost {
		t.Errorf("method = %q, want %q", capturedRequest.Method, http.MethodPost)
	}
	if capturedRequest.URL.Path != "/b" {
		t.Errorf("path = %q, want %q", capturedRequest.URL.Path, "/b")
	}
	if got := capturedRequest.Header.Get(masterKeyHeader); got != "test-key" {
		t.Errorf("%s = %q, want %q", masterKeyHeader, got, "test-key")
	}
	if got := capturedRequest.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if got := capturedRequest.Header.Get(binNameHeader); got != "my-bin" {
		t.Errorf("%s = %q, want my-bin", binNameHeader, got)
	}
	if got := capturedRequest.Header.Get(binPrivateHeader); got != "true" {
		t.Errorf("%s = %q, want true", binPrivateHeader, got)
	}
	if capturedBodyError != nil {
		t.Fatalf("read request body: %v", capturedBodyError)
	}
	if got, expectedBody := string(capturedBody), `{"hello":"world"}`; got != expectedBody {
		t.Errorf("body = %q, want %q", got, expectedBody)
	}
	if created != expected {
		t.Fatalf("Create() = %+v, want %+v", created, expected)
	}
}

func TestGetUsesBinRoute(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	expected := `{"value":2}`
	var capturedRequest *http.Request
	client := testClient(httpClientFunc(func(request *http.Request) (*http.Response, error) {
		capturedRequest = request
		return jsonResponse(http.StatusOK, `{"record":{"value":2},"metadata":{"id":"bin-id"}}`), nil
	}))

	// Act - выполняем функцию
	record, err := client.Get(context.Background(), "bin-id")

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if capturedRequest == nil {
		t.Fatal("Get() did not send an HTTP request")
	}
	if got := capturedRequest.Method; got != http.MethodGet {
		t.Errorf("method = %q, want %q", got, http.MethodGet)
	}
	if got := capturedRequest.URL.Path; got != "/b/bin-id" {
		t.Errorf("path = %q, want %q", got, "/b/bin-id")
	}
	if got := string(record); got != expected {
		t.Fatalf("Get() = %s, want record %s", got, expected)
	}
}

func TestUpdateUsesBinRoute(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	expected := `{"record":{"value":2},"metadata":{"id":"bin-id"}}`
	var capturedRequest *http.Request
	var capturedBody []byte
	var capturedBodyError error
	client := testClient(httpClientFunc(func(request *http.Request) (*http.Response, error) {
		capturedRequest = request
		capturedBody, capturedBodyError = io.ReadAll(request.Body)
		return jsonResponse(http.StatusOK, `{"record":{"value":2},"metadata":{"id":"bin-id"}}`), nil
	}))

	// Act - выполняем функцию
	response, err := client.Update(context.Background(), "bin-id", []byte(`{"value":2}`))

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if capturedRequest == nil {
		t.Fatal("Update() did not send an HTTP request")
	}
	if got := capturedRequest.Method; got != http.MethodPut {
		t.Errorf("method = %q, want %q", got, http.MethodPut)
	}
	if got := capturedRequest.URL.Path; got != "/b/bin-id" {
		t.Errorf("path = %q, want %q", got, "/b/bin-id")
	}
	if got := capturedRequest.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if capturedBodyError != nil {
		t.Fatalf("read request body: %v", capturedBodyError)
	}
	if got, expectedBody := string(capturedBody), `{"value":2}`; got != expectedBody {
		t.Errorf("body = %q, want %q", got, expectedBody)
	}
	if got := string(response); got != expected {
		t.Fatalf("Update() = %s, want response %s", got, expected)
	}
}

func TestDeleteUsesBinRoute(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	var capturedRequest *http.Request
	client := testClient(httpClientFunc(func(request *http.Request) (*http.Response, error) {
		capturedRequest = request
		return jsonResponse(http.StatusNoContent, ""), nil
	}))

	// Act - выполняем функцию
	err := client.Delete(context.Background(), "bin-id")

	// Assert - проверка результата с expected
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if capturedRequest == nil {
		t.Fatal("Delete() did not send an HTTP request")
	}
	if got := capturedRequest.Method; got != http.MethodDelete {
		t.Errorf("method = %q, want %q", got, http.MethodDelete)
	}
	if got := capturedRequest.URL.Path; got != "/b/bin-id" {
		t.Errorf("path = %q, want %q", got, "/b/bin-id")
	}
}

func TestGetRejectsResponseWithoutRecord(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	expectedError := "does not contain a record"
	client := testClient(httpClientFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{"metadata":{"id":"bin-id"}}`), nil
	}))

	// Act - выполняем функцию
	_, err := client.Get(context.Background(), "bin-id")

	// Assert - проверка результата с expected
	if err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("Get() error = %v, want missing record error", err)
	}
}

func TestGetRejectsMalformedResponse(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	expectedError := "decode response"
	client := testClient(httpClientFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{"record":`), nil
	}))

	// Act - выполняем функцию
	_, err := client.Get(context.Background(), "bin-id")

	// Assert - проверка результата с expected
	if err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("Get() error = %v, want decode response error", err)
	}
}

func TestAPIErrorIncludesStatusAndMessage(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	expectedStatus := http.StatusUnauthorized
	expectedMessage := "invalid key"
	client := testClient(httpClientFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(expectedStatus, `{"message":"invalid key"}`), nil
	}))

	// Act - выполняем функцию
	_, err := client.Get(context.Background(), "bin-id")

	// Assert - проверка результата с expected
	var apiError *APIError
	if !errors.As(err, &apiError) {
		t.Fatalf("Get() error = %v, want *APIError", err)
	}
	if apiError.StatusCode != expectedStatus || apiError.Message != expectedMessage {
		t.Fatalf("APIError = %+v", apiError)
	}
}

func TestCreateRejectsInvalidJSONBeforeRequest(t *testing.T) {
	// Arrange - подготовка, expected результат, данные для функции
	expectedError := ErrInvalidJSON
	client := testClient(httpClientFunc(func(*http.Request) (*http.Response, error) {
		panic("HTTP request must not be sent")
	}))

	// Act - выполняем функцию
	_, err := client.Create(context.Background(), []byte(`{"broken"`), "my-bin")

	// Assert - проверка результата с expected
	if !errors.Is(err, expectedError) {
		t.Fatalf("Create() error = %v, want %v", err, expectedError)
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
