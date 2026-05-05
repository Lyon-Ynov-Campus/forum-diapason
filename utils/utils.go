package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

const SessionCookieName = "session_id"
const SessionDuration = 24 * time.Hour

//Session 

// GenerateSessionID génère un token aléatoire sécurisé (64 caractères hex)
func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// SetSessionCookie pose le cookie de session dans la réponse HTTP
func SetSessionCookie(w http.ResponseWriter, sessionID string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true, // protège contre XSS (JS ne peut pas lire le cookie)
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // décommenter en production (HTTPS)
	})
}

// ClearSessionCookie supprime le cookie côté client
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// GetUserIDFromSession lit le cookie, vérifie la session en base et retourne l'user_id
// Retourne 0 si non connecté ou session expirée
func GetUserIDFromSession(r *http.Request, db *sql.DB) int {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return 0
	}

	var userID int
	query := `SELECT user_id FROM sessions WHERE id = ? AND expires_at > datetime('now')`
	err = db.QueryRow(query, cookie.Value).Scan(&userID)
	if err != nil {
		return 0
	}
	return userID
}

// RequireAuth vérifie que l'utilisateur est connecté.
// Retourne l'userID si OK, sinon écrit 401 et retourne 0.
func RequireAuth(w http.ResponseWriter, r *http.Request, db *sql.DB) int {
	userID := GetUserIDFromSession(r, db)
	if userID == 0 {
		http.Error(w, `{"error":"non authentifié"}`, http.StatusUnauthorized)
		return 0
	}
	return userID
}

//Validation 

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func IsValidPassword(password string) bool {
	return len(password) >= 8
}

func IsValidPseudo(pseudo string) bool {
	return len(pseudo) >= 3 && len(pseudo) <= 30
}