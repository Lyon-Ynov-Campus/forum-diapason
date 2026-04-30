package handlers

import (
	"encoding/json"
	"net/http"

	"forum-diapason/internal/usecases"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUseCase usecases.UserUseCase
}

type UpdateProfileRequest struct {
	Email    string `json:"email"`
	Pseudo   string `json:"pseudo"`
	PhotoURL string `json:"photo_url"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func NewUserHandler(userUseCase usecases.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	user, err := h.userUseCase.GetProfile(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.userUseCase.UpdateProfile(userID, req.Email, req.Pseudo, req.PhotoURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Profile updated successfully"}`))
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.userUseCase.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Password changed successfully"}`))
}

func (h *UserHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userIDStr, exists := vars["userId"]
	if !exists {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	_ = userIDStr // TODO: Implement user posts retrieval

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	_ = limitStr  // TODO: Use limit for pagination
	_ = offsetStr // TODO: Use offset for pagination

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "User posts endpoint - implementation pending"}`))
}
