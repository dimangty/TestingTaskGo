package api

import (
	"context"
	"net/http"
	"testing"

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
