package handlers

// Pages d'auth : inscription, connexion, deco
// Tout fonctionne en rendu HTML classique : form POST → redirect en cas de
// succes, ou re-affichage de la page avec un message d'erreur

import (
	"net/http"

	"forum-diapason/services"
)

// LoginPage gere le GET (affichage du form) et le POST (tentative de co)
// Un user deja co est renvoye chez lui sans meme voir le formulaire
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if RedirectIfAuthed(w, r) {
		return
	}
	if r.Method == http.MethodGet {
		RenderPage(w, r, "login", nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	emailOrPseudo := r.FormValue("email_or_pseudo")
	password := r.FormValue("password")

	// On garde la valeur saisie pour la re-injecter dans le form en cas d'echec
	form := map[string]any{"EmailOrPseudo": emailOrPseudo}
	renderErr := func(msg string) {
		form["Error"] = msg
		RenderPage(w, r, "login", form)
	}

	user, err := services.Login(db, emailOrPseudo, password)
	if err != nil {
		renderErr(err.Error())
		return
	}
	if err := services.CreateSession(db, w, user.ID); err != nil {
		renderErr("erreur de session")
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// RegisterPage : meme decoupage que LoginPage, avec en plus la confirmation
// du mdp (verifiee ici, pas dans le service)
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if RedirectIfAuthed(w, r) {
		return
	}
	if r.Method == http.MethodGet {
		RenderPage(w, r, "register", nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	nom := r.FormValue("nom")
	pseudo := r.FormValue("pseudo")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm := r.FormValue("password_confirm")

	form := map[string]any{
		"Nom":    nom,
		"Pseudo": pseudo,
		"Email":  email,
	}
	renderErr := func(msg string) {
		form["Error"] = msg
		RenderPage(w, r, "register", form)
	}

	if password != confirm {
		renderErr("les mots de passe ne correspondent pas")
		return
	}
	user, err := services.Register(db, nom, pseudo, email, password)
	if err != nil {
		renderErr(err.Error())
		return
	}
	if err := services.CreateSession(db, w, user.ID); err != nil {
		renderErr("erreur de session")
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutPage est volontairement POST-only : sans ca, un simple <img src="/logout">
// pose sur un site tiers suffirait a deco le user (CSRF)
func LogoutPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	services.Logout(db, w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Me renvoie le user co en JSON
// Pas utilise par le rendu des pages, garde sous la main pour un eventuel
// appel AJAX plus tard
func Me(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		sendError(w, http.StatusUnauthorized, "non authentifié")
		return
	}
	sendJSON(w, http.StatusOK, user)
}
