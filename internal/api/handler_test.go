package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomePageIntroducesInteractiveLibrary(t *testing.T) {
	response := request(t, NewHandler(), http.MethodGet, "/", "")
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Content-Type"); !strings.Contains(got, "text/html") {
		t.Fatalf("content type = %q, want text/html", got)
	}
	for _, text := range []string{"Biblioteca de pruebas", "Agregar libro", "Buscar título o autor", "/api/books"} {
		if !strings.Contains(response.Body.String(), text) {
			t.Errorf("home page does not contain %q", text)
		}
	}
}

func TestHealth(t *testing.T) {
	response := request(t, NewHandler(), http.MethodGet, "/health", "")
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status body = %q, want ok", body["status"])
	}
}

func TestCreateBookRequiresTitleAndAuthor(t *testing.T) {
	response := request(t, NewHandler(), http.MethodPost, "/api/books", `{"title":"A book"}`)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestGetBookReturnsNotFound(t *testing.T) {
	response := request(t, NewHandler(), http.MethodGet, "/api/books/999", "")
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}

func request(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	handler.ServeHTTP(recorder, req)
	return recorder
}
