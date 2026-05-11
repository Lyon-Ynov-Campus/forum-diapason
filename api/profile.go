package main

import (
	"encoding/json"
	"net/http"

	"forum-diapason/services"
	"forum-diapason/utils"
)

func profileUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	userID := utils.GetUserIDFromSession(r, db)
	if userID == 0 {
		writeError(w, http.StatusUnauthorized, "non authentifié")
		return
	}
	var body struct {
		Nom    string `json:"nom"`
		Pseudo string `json:"pseudo"`
		Email  string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	if err := services.UpdateProfile(db, userID, body.Nom, body.Pseudo, body.Email); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "profil mis à jour"})
}

func profilePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	userID := utils.GetUserIDFromSession(r, db)
	if userID == 0 {
		writeError(w, http.StatusUnauthorized, "non authentifié")
		return
	}
	var body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	if err := services.UpdatePassword(db, userID, body.OldPassword, body.NewPassword); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "mot de passe mis à jour"})
}
