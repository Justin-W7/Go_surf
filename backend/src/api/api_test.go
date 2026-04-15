package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchURL_Success(t *testing.T) {
	// mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	data, err := fetchURL(server.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(data) != `{"ok":true}` {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestFetchURL_EmptyURL(t *testing.T) {
	_, err := fetchURL("")

	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestFetchURL_InvalidURL(t *testing.T) {
	_, err := fetchURL("htp://bad-url")

	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
