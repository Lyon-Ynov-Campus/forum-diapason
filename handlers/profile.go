package handlers

// Pages du profil : affichage + modification (infos, avatar, mdp)

import (
	"fmt"
	"net/http"

	"forum-diapason/models"
	"forum-diapason/services"
)

// GET /profile       → profil du user connecté
// GET /profile?id=2  → profil d'un autre utilisateur
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	if RequirePageAuth(w, r) == 0 {
		return
	}

	var profileUser *models.User

	if idStr := r.URL.Query().Get("id"); idStr != "" {
		var userID int
		if _, err := fmt.Sscanf(idStr, "%d", &userID); err == nil {
			profileUser, _ = services.GetUserByID(db, userID)
		}
	}

	// Pas d'id ou user introuvable → on affiche le user connecté
	if profileUser == nil {
		profileUser = currentUser(r)
	}

	loggedIn := currentUser(r)
	isOwn := loggedIn != nil && profileUser != nil && loggedIn.ID == profileUser.ID
	RenderPage(w, r, "profile", map[string]any{
		"ProfileUser": profileUser,
		"IsOwnProfile": isOwn,
	})
}

// GET  /profile/edit  → form pre-rempli avec les valeurs actuelles
// POST /profile/edit  → MAJ des infos (nom, pseudo, email)
func ProfileEditPage(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		// Les cles Nom/Pseudo/... (vs .User.Nom) servent a pouvoir re-injecter
		// les valeurs *saisies* en cas d'erreur, sans reecraser avec les valeurs en base
		RenderPage(w, r, "profile-edit", map[string]any{
			"Nom":    user.Nom,
			"Pseudo": user.Pseudo,
			"Email":  user.Email,
		})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	nom := r.FormValue("nom")
	pseudo := r.FormValue("pseudo")
	email := r.FormValue("email")

	form := map[string]any{
		"Nom":    nom,
		"Pseudo": pseudo,
		"Email":  email,
	}

	if err := services.UpdateProfile(db, user.ID, nom, pseudo, email); err != nil {
		form["Error"] = err.Error()
		RenderPage(w, r, "profile-edit", form)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/profile?id=%d", user.ID), http.StatusSeeOther)
}

// POST /profile/avatar  → upload d'une nouvelle photo de profil (multipart)
func ProfileAvatarPage(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// MaxBytesReader avorte la requete si le body depasse la limite, avant
	// meme que ParseMultipartForm n'essaie de tout charger en memoire/disque
	r.Body = http.MaxBytesReader(w, r.Body, services.AvatarMaxSize)
	if err := r.ParseMultipartForm(services.AvatarMaxSize); err != nil {
		rerenderAvatarErr(w, r, user, "image trop lourde (max 2 Mo)")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		rerenderAvatarErr(w, r, user, "fichier manquant")
		return
	}
	defer file.Close()

	if _, err := services.UpdateAvatar(db, user.ID, file, header); err != nil {
		rerenderAvatarErr(w, r, user, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/profile?id=%d", user.ID), http.StatusSeeOther)
}

// POST /profile/avatar/delete  → supprime la photo de profil (fichier + DB)
func ProfileAvatarDeletePage(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if err := services.DeleteAvatar(db, user.ID); err != nil {
		rerenderAvatarErr(w, r, user, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/profile?id=%d", user.ID), http.StatusSeeOther)
}

// POST /profile/password  → change le mdp (ancien + nouveau + confirmation)
func ProfilePasswordPage(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	oldPassword := r.FormValue("old_password")
	newPassword := r.FormValue("new_password")
	confirm := r.FormValue("new_password_confirm")

	rerender := func(errMsg string) {
		RenderPage(w, r, "profile-edit", map[string]any{
			"Nom":           user.Nom,
			"Pseudo":        user.Pseudo,
			"Email":         user.Email,
			"PasswordError": errMsg,
		})
	}

	if newPassword != confirm {
		rerender("les nouveaux mots de passe ne correspondent pas")
		return
	}
	if err := services.UpdatePassword(db, user.ID, oldPassword, newPassword); err != nil {
		rerender(err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/profile?id=%d", user.ID), http.StatusSeeOther)
}

// rerenderAvatarErr : helper pour rerender /profile/edit avec une erreur cote
// upload, sans perdre les valeurs actuelles dans les autres sections du form
func rerenderAvatarErr(w http.ResponseWriter, r *http.Request, user *models.User, msg string) {
	RenderPage(w, r, "profile-edit", map[string]any{
		"Nom":         user.Nom,
		"Pseudo":      user.Pseudo,
		"Email":       user.Email,
		"AvatarError": msg,
	})
}
