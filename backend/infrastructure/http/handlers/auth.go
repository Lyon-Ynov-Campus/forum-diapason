package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"forum-diapason/pkg/config"
	"forum-diapason/internal/usecases"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userUseCase usecases.UserUseCase
	config      *config.Config
}

type LoginRequest struct {
	EmailOrPseudo string `json:"email_or_pseudo"`
	Password      string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Pseudo   string `json:"pseudo"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

func NewAuthHandler(userUseCase usecases.UserUseCase, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userUseCase: userUseCase,
		config:      cfg,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.Register(req.Email, req.Pseudo, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.generateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		User:  user,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.Login(req.EmailOrPseudo, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.generateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		User:  user,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) generateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(h.config.JWT.Secret))
}
