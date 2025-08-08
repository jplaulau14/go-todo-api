package todo

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type HTTPHandler struct {
	repo Repository
}

func NewHTTPHandler(repo Repository) *HTTPHandler {
	return &HTTPHandler{repo: repo}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	// Support both /todos and /todos/
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/todos" { // let other paths fall through to /todos/
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.list(w, r)
			return
		case http.MethodPost:
			h.create(w, r)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		// /todos/ -> list, create
		// /todos/{id} -> get, update, delete
		path := strings.TrimPrefix(r.URL.Path, "/todos/")
		if path == "" { // exactly /todos/
			switch r.Method {
			case http.MethodGet:
				h.list(w, r)
				return
			case http.MethodPost:
				h.create(w, r)
				return
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
		}

		// path contains id (could have trailing slash removed already)
		id := strings.TrimSuffix(path, "/")
		switch r.Method {
		case http.MethodGet:
			h.get(w, r, id)
			return
		case http.MethodPatch:
			h.update(w, r, id)
			return
		case http.MethodDelete:
			h.delete(w, r, id)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *HTTPHandler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	t, err := h.repo.Create(r.Context(), req.Title)
	if err != nil {
		http.Error(w, "could not create", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *HTTPHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "could not list", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *HTTPHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	t, err := h.repo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "could not get", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *HTTPHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	updated, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "could not update", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *HTTPHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "could not delete", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
