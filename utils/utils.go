package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"net/mail"
	"regexp"
	"time"
)

const SessionCookieName = "session_id"
const SessionDuration = 24 * time.Hour

// CookieSecure pose le flag Secure sur le cookie de session
// A mettre a true en HTTPS (prod) et false en HTTP (dev)
// Configure depuis main.go via la variable d'env COOKIE_SECURE
var CookieSecure bool

// --- Session ---

// GenerateSessionID renvoie 32 octets aleatoires encodes en hex (64 caracteres)
// Suffisant pour un id de session non devinable
func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// SetSessionCookie pose le cookie de session sur la reponse
//   - HttpOnly : empeche le JS de lire le cookie (mitigation XSS)
//   - SameSite=Lax : bloque la majorite des requetes cross-site (mitigation CSRF)
func SetSessionCookie(w http.ResponseWriter, sessionID string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   CookieSecure,
	})
}

// ClearSessionCookie demande au browser d'oublier le cookie
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// GetUserIDFromSession lit le cookie, verifie la session en base et renvoie
// l'user_id
// Retourne 0 si pas de cookie, session inconnue ou expiree
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

// RequireAuth est la garde cote API : si non co, on repond 401 JSON et
// retourne 0
// Pour les pages HTML, voir handlers.RequirePageAuth
func RequireAuth(w http.ResponseWriter, r *http.Request, db *sql.DB) int {
	userID := GetUserIDFromSession(r, db)
	if userID == 0 {
		http.Error(w, `{"error":"non authentifié"}`, http.StatusUnauthorized)
		return 0
	}
	return userID
}

// --- Validation ---

var pseudoRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// IsValidEmail accepte une adresse seule et refuse la forme "Nom <email>"
// Par defaut mail.ParseAddress autorise les deux, on compare donc avec
// l'adresse parsee pour etre strict
func IsValidEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return addr.Address == email
}

func IsValidPassword(password string) bool {
	return len(password) >= 8
}

func IsValidPseudo(pseudo string) bool {
	if len(pseudo) < 3 || len(pseudo) > 30 {
		return false
	}
	return pseudoRegex.MatchString(pseudo)
}

func IsValidNom(nom string) bool {
	return len(nom) >= 1 && len(nom) <= 50
}
