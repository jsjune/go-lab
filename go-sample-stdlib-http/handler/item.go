package handler

import (
	"encoding/json"
	"errors"
	"go-sample-stdlib-http/store"
	"net/http"
	"strconv"
	"strings"
)

type Store interface {
	List() []store.Item
	Get(id int64) (*store.Item, error)
	Create(name, description string) store.Item
	Update(id int64, name, description string) (*store.Item, error)
	Delete(id int64) error
}

type ItemHandler struct {
	store Store
}

func NewItemHandler(s Store) *ItemHandler {
	return &ItemHandler{store: s}
}

// RegisterRoutes — Gin 없이 표준 ServeMux로 라우팅
// Gin이 내부에서 하는 일을 직접 구현
func (h *ItemHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/items", h.handleItems)    // GET, POST
	mux.HandleFunc("/items/", h.handleItem)    // GET, PUT, DELETE /items/:id
}

// /items → 메서드에 따라 분기
func (h *ItemHandler) handleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// /items/:id → 메서드에 따라 분기
func (h *ItemHandler) handleItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.get(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *ItemHandler) list(w http.ResponseWriter, r *http.Request) {
	items := h.store.List()
	writeJSON(w, http.StatusOK, items)
}

func (h *ItemHandler) get(w http.ResponseWriter, r *http.Request, id int64) {
	item, err := h.store.Get(id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	item := h.store.Create(req.Name, req.Description)
	writeJSON(w, http.StatusCreated, item)
}

func (h *ItemHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	item, err := h.store.Update(id, req.Name, req.Description)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) delete(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.store.Delete(id); errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// /items/1 → id=1 추출
func parseID(r *http.Request) (int64, error) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/items/"), "/")
	return strconv.ParseInt(parts[0], 10, 64)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
