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
