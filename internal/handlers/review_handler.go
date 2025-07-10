package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"reviews/internal/models"
	"reviews/internal/services"
)

type ReviewHandler struct {
	Service *services.ReviewService
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	var rev models.Reviews
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	created, err := h.Service.Create(r.Context(), rev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

func (h *ReviewHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.Service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	rev, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rev)
}

func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	var rev models.Reviews
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	rev.ID = id
	updated, err := h.Service.Update(r.Context(), rev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	if err := h.Service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
