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
	req = httptest.NewRequest(http.MethodGet, "/todos/", nil)
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
