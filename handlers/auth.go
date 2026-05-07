package handlers

// Inscription, connexion, déconnexion

import (
	"encoding/json"
	"forum-diapason/services"
	"forum-diapason/utils"
	"net/http"
)

// POST /api/auth/register
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	var body struct {
		Nom      string `json:"nom"`
		Pseudo   string `json:"pseudo"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	user, err := services.Register(db, body.Nom, body.Pseudo, body.Email, body.Password)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := services.CreateSession(db, w, user.ID); err != nil {
		sendError(w, http.StatusInternalServerError, "erreur de session")
		return
	}

	sendJSON(w, http.StatusCreated, map[string]any{
		"message": "inscription réussie",
		"user":    user,
	})
}

// POST /api/auth/login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	var body struct {
		EmailOrPseudo string `json:"email_or_pseudo"`
		Password      string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	user, err := services.Login(db, body.EmailOrPseudo, body.Password)
	if err != nil {
		sendError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := services.CreateSession(db, w, user.ID); err != nil {
		sendError(w, http.StatusInternalServerError, "erreur de session")
		return
	}

	sendJSON(w, http.StatusOK, map[string]any{
		"message": "connexion réussie",
		"user":    user,
	})
}

// POST /api/auth/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	services.Logout(db, w, r)
	sendJSON(w, http.StatusOK, map[string]string{"message": "déconnexion réussie"})
}

// GET /api/auth/me → retourne l'utilisateur connecté
// DELETE /api/auth/me → supprime le compte
func Me(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		userID := utils.RequireAuth(w, r, db)
		if userID == 0 {
			return
		}
		var user struct {
			ID       int    `json:"id"`
			Nom      string `json:"nom"`
			Pseudo   string `json:"pseudo"`
			Email    string `json:"email"`
			PhotoURL string `json:"photo_url"`
		}
		err := db.QueryRow(
			`SELECT id, nom, pseudo, email, photo_url FROM users WHERE id = ?`, userID,
		).Scan(&user.ID, &user.Nom, &user.Pseudo, &user.Email, &user.PhotoURL)
		if err != nil {
			sendError(w, http.StatusNotFound, "utilisateur introuvable")
			return
		}
		sendJSON(w, http.StatusOK, user)

	case http.MethodDelete:
		userID := utils.RequireAuth(w, r, db)
		if userID == 0 {
			return
		}
		// Déconnexion d'abord, puis suppression
		services.Logout(db, w, r)
		if _, err := db.Exec(`DELETE FROM users WHERE id = ?`, userID); err != nil {
			sendError(w, http.StatusInternalServerError, "erreur suppression du compte")
			return
		}
		sendJSON(w, http.StatusOK, map[string]string{"message": "compte supprimé"})
	}
}
