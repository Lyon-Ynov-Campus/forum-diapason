package main

import (
	"encoding/json"
	"net/http"

	"forum-diapason/services"
	"forum-diapason/utils"
)

func authRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	var body struct {
		Nom      string `json:"nom"`
		Pseudo   string `json:"pseudo"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	user, err := services.Register(db, body.Nom, body.Pseudo, body.Email, body.Password)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := services.CreateSession(db, w, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "erreur session")
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func authLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	var body struct {
		EmailOrPseudo string `json:"email_or_pseudo"`
		Password      string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	user, err := services.Login(db, body.EmailOrPseudo, body.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err := services.CreateSession(db, w, user.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "erreur session")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func authLogout(w http.ResponseWriter, r *http.Request) {
	services.Logout(db, w, r)
	writeJSON(w, http.StatusOK, map[string]string{"message": "déconnecté"})
}

func authMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	userID := utils.GetUserIDFromSession(r, db)
	if userID == 0 {
		writeError(w, http.StatusUnauthorized, "non authentifié")
		return
	}
	user, err := services.GetUserByID(db, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func authForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	token, _, err := services.CreatePasswordReset(db, body.Email)
	if err == nil {
		front := getEnv("FRONTEND_ORIGIN", "http://localhost:8080")
		link  := front + "/reset-password?token=" + token
		utils.SendMail(body.Email, "Réinitialisation de votre mot de passe Diapason",
			"Bonjour,\n\nClique sur ce lien pour réinitialiser ton mot de passe :\n"+link+"\n\nCe lien expire dans 1 heure.")
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "si l'email existe, un lien a été envoyé"})
}

func authResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	var body struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	if err := services.ResetPasswordWithToken(db, body.Token, body.Password); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "mot de passe réinitialisé"})
}
