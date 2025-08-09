package todo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupServer() http.Handler {
	repo := NewInMemoryRepository()
	h := NewHTTPHandler(repo)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func TestHTTP_CRUD(t *testing.T) {
	srv := setupServer()

	// Create
	body := bytes.NewBufferString(`{"title":"task1"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status: %d", w.Code)
	}
	var created Todo
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	// Get
	req = httptest.NewRequest(http.MethodGet, "/todos/"+created.ID, nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get status: %d", w.Code)
	}

	// List
	req = httptest.NewRequest(http.MethodGet, "/todos/?limit=5&offset=0", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status: %d", w.Code)
	}
	b, _ := io.ReadAll(w.Body)
	if !bytes.Contains(b, []byte(created.ID)) {
		t.Fatalf("list missing created id")
	}

	// Update
	updateBody := bytes.NewBufferString(`{"completed":true}`)
	req = httptest.NewRequest(http.MethodPatch, "/todos/"+created.ID, updateBody)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("update status: %d", w.Code)
	}

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/todos/"+created.ID, nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete status: %d", w.Code)
	}
}

func TestHTTP_ListPaginationParams(t *testing.T) {
	srv := setupServer()

	// Seed some todos
	for i := 0; i < 5; i++ {
		body := bytes.NewBufferString(`{"title":"t"}`)
		req := httptest.NewRequest(http.MethodPost, "/todos/", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("seed create status: %d", w.Code)
		}
	}

	// limit=2 offset=1
	req := httptest.NewRequest(http.MethodGet, "/todos/?limit=2&offset=1", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status: %d", w.Code)
	}
	var items []Todo
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	// Invalid/negative params should be clamped to defaults
	req = httptest.NewRequest(http.MethodGet, "/todos/?limit=-1&offset=-5", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status: %d", w.Code)
	}
	items = nil
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("expected some items with clamped params")
	}
}

func TestHTTP_CreateWithoutTrailingSlash(t *testing.T) {
	srv := setupServer()
	body := bytes.NewBufferString(`{"title":"task2"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status without trailing slash: %d", w.Code)
	}
}

func TestHTTP_Validation_ContentType(t *testing.T) {
	srv := setupServer()
	body := bytes.NewBufferString(`{"title":"task"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	// Missing Content-Type
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415, got %d", w.Code)
	}
}

func TestHTTP_Validation_UnknownFields(t *testing.T) {
	srv := setupServer()
	body := bytes.NewBufferString(`{"title":"t","unknown":true}`)
	req := httptest.NewRequest(http.MethodPost, "/todos/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
