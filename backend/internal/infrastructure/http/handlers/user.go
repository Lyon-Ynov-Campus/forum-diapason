package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// UserHandlers handles user-related routes
type UserHandlers struct {
	db *sql.DB
}

// NewUserHandlers creates a new UserHandlers instance
func NewUserHandlers(db *sql.DB) *UserHandlers {
	return &UserHandlers{db: db}
}

// GetProfile retrieves a user profile
func (h *UserHandlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userID := vars["id"]

	// TODO: Implement get profile logic
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": userID})
}

// UpdateProfile updates a user profile
func (h *UserHandlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userID := vars["id"]

	var req struct {
		Email    string `json:"email"`
		Pseudo   string `json:"pseudo"`
		PhotoURL string `json:"photo_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// TODO: Implement update profile logic
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"user_id": userID, "message": "Profile updated"})
}

// SearchUsers searches for users by pseudo
func (h *UserHandlers) SearchUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")

	// TODO: Implement search users logic
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"query": query, "results": []interface{}{}})
}
