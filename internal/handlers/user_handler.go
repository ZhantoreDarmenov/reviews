package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"reviews/internal/models"
	"reviews/internal/services"
)

type UserHandler struct {
	Service *services.UserService
}

// SignUp registers a new user and returns auth tokens.
func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.Service.SignUp(r.Context(), user)
	if err != nil {
		log.Printf("sign up error: %v", err)
		http.Error(w, "failed to sign up", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.Service.SignIn(r.Context(), req.Login, req.Password)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "user id missing", http.StatusUnauthorized)
		return
	}
	if err := h.Service.Logout(r.Context(), userID); err != nil {
		log.Printf("logout error: %v", err)
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "logout successful"})
}
