package todo

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jplaulau14/go-todo-api/internal/reqctx"
)

type HTTPHandler struct {
	repo   Repository
	logger *slog.Logger
}

func NewHTTPHandler(repo Repository) *HTTPHandler {
	return &HTTPHandler{repo: repo, logger: slog.Default()}
}

func (h *HTTPHandler) WithLogger(logger *slog.Logger) *HTTPHandler {
	if logger == nil {
		return h
	}
	h.logger = logger
	return h
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/todos" {
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
			writeError(w, r, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	})
	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/todos/")
		if path == "" {
			switch r.Method {
			case http.MethodGet:
				h.list(w, r)
				return
			case http.MethodPost:
				h.create(w, r)
				return
			default:
				writeError(w, r, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
		}

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
			writeError(w, r, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type errorResponse struct {
	Error     string `json:"error"`
	Status    int    `json:"status"`
	RequestID string `json:"request_id,omitempty"`
}

func writeError(w http.ResponseWriter, r *http.Request, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message, Status: status, RequestID: reqctx.GetRequestID(r.Context())})
}

func isJSON(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	if ct == "" {
		return false
	}
	// Accept application/json and application/json; charset=UTF-8
	ct = strings.ToLower(ct)
	return strings.HasPrefix(ct, "application/json")
}

func (h *HTTPHandler) create(w http.ResponseWriter, r *http.Request) {
	if !isJSON(r) {
		writeError(w, r, http.StatusUnsupportedMediaType, "content-type must be application/json")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req CreateTodoRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			writeError(w, r, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		h.logger.Warn("invalid json", "error", err, "request_id", reqctx.GetRequestID(r.Context()))
		writeError(w, r, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Title == "" {
		writeError(w, r, http.StatusBadRequest, "title is required")
		return
	}
	t, err := h.repo.Create(r.Context(), req.Title)
	if err != nil {
		h.logger.Error("could not create todo", "error", err, "request_id", reqctx.GetRequestID(r.Context()))
		writeError(w, r, http.StatusInternalServerError, "could not create")
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *HTTPHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		h.logger.Error("could not list todos", "error", err, "request_id", reqctx.GetRequestID(r.Context()))
		writeError(w, r, http.StatusInternalServerError, "could not list")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *HTTPHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	t, err := h.repo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, r, http.StatusNotFound, "todo not found")
			return
		}
		h.logger.Error("could not get todo", "id", id, "error", err, "request_id", reqctx.GetRequestID(r.Context()))
		writeError(w, r, http.StatusInternalServerError, "could not get")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *HTTPHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	if !isJSON(r) {
		writeError(w, r, http.StatusUnsupportedMediaType, "content-type must be application/json")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req UpdateTodoRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			writeError(w, r, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		h.logger.Warn("invalid json", "error", err, "request_id", reqctx.GetRequestID(r.Context()))
		writeError(w, r, http.StatusBadRequest, "invalid json")
		return
	}
	updated, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, r, http.StatusNotFound, "todo not found")
			return
		}
		h.logger.Error("could not update todo", "id", id, "error", err, "request_id", reqctx.GetRequestID(r.Context()))
		writeError(w, r, http.StatusInternalServerError, "could not update")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *HTTPHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, r, http.StatusNotFound, "todo not found")
			return
		}
		h.logger.Error("could not delete todo", "id", id, "error", err, "request_id", reqctx.GetRequestID(r.Context()))
		writeError(w, r, http.StatusInternalServerError, "could not delete")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
