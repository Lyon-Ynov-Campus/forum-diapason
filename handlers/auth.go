package handlers

// Pages d'auth : inscription, connexion, déconnexion.
// Tout fonctionne en rendu HTML classique : form POST → redirect en cas de
// succès, ou ré-affichage de la page avec un message d'erreur.

import (
	"net/http"

	"forum-diapason/services"
)

// LoginPage gère le GET (affichage du form) et le POST (tentative de connexion).
// Un user déjà connecté est renvoyé chez lui sans même voir le formulaire.
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

	// On garde la valeur saisie pour la ré-injecter dans le form en cas d'échec.
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

// RegisterPage : même découpage que LoginPage, avec en plus la confirmation
// du mot de passe (vérifiée ici, pas dans le service).
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

// LogoutPage est volontairement POST-only : sans ça, un simple <img src="/logout">
// posé sur un site tiers suffirait à déconnecter l'utilisateur (CSRF).
func LogoutPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	services.Logout(db, w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Me renvoie l'utilisateur connecté en JSON. Pas utilisé par le rendu des pages,
// gardé sous la main pour un éventuel appel AJAX plus tard.
func Me(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		sendError(w, http.StatusUnauthorized, "non authentifié")
		return
	}
	sendJSON(w, http.StatusOK, user)
}
