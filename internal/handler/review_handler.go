package handler

import (
	"encoding/json"
	"net/http"
	"reviews/internal/models"
	"reviews/internal/service"
	"strconv"
	"strings"
)

type ReviewHandler struct {
	svc  *service.ReviewService
	auth *service.AuthService
}

func NewReviewHandler(s *service.ReviewService, a *service.AuthService) *ReviewHandler {
	return &ReviewHandler{svc: s, auth: a}
}

func (h *ReviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/reviews" || r.URL.Path == "/reviews/" {
		h.handleCollection(w, r)
	} else {
		h.handleItem(w, r)
	}
}

func (h *ReviewHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		reviews, _ := h.svc.GetAll()
		json.NewEncoder(w).Encode(reviews)
	case http.MethodPost:
		if !h.authorized(r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var rev models.Review
		if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		created, _ := h.svc.Create(rev)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ReviewHandler) handleItem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/reviews/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		rev, ok, _ := h.svc.GetByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(rev)
	case http.MethodPut:
		if !h.authorized(r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var rev models.Review
		if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		rev.ID = id
		updated, ok, _ := h.svc.Update(rev)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(updated)
	case http.MethodDelete:
		if !h.authorized(r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		deleted := h.svc.Delete(id)
		if !deleted {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ReviewHandler) authorized(r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if token == "" {
		return false
	}
	return h.auth.Authenticate(token)
}
