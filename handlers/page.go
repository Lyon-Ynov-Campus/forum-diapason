package handlers

// Helpers partagés par tous les handlers de pages : rendu de template
// + gardes d'auth pensées pour le HTML (redirections plutôt que JSON).

import (
	"html/template"
	"net/http"

	"forum-diapason/models"
	"forum-diapason/utils"
)

// RenderPage charge la page demandée + les composants partagés (header, footer…)
// et l'exécute. On glisse au passage l'utilisateur connecté sous la clé "User"
// pour que le header sache afficher le pseudo (ou les boutons sign in/register).
func RenderPage(w http.ResponseWriter, r *http.Request, name string, data map[string]any) {
	if data == nil {
		data = map[string]any{}
	}
	if _, ok := data["User"]; !ok {
		data["User"] = currentUser(r)
	}
	t := template.Must(template.ParseGlob("./frontend/components/*.html"))
	t = template.Must(t.ParseFiles("./frontend/pages/" + name + ".html"))
	t.ExecuteTemplate(w, name+".html", data)
}

// RequirePageAuth est la version "page" de utils.RequireAuth : si l'user n'est
// pas connecté, on redirige vers /login plutôt que de répondre 401 en JSON.
// Retourne 0 si la redirection a été émise, sinon l'userID.
func RequirePageAuth(w http.ResponseWriter, r *http.Request) int {
	userID := utils.GetUserIDFromSession(r, db)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 0
	}
	return userID
}

// RedirectIfAuthed sert sur /login et /register : si on est déjà connecté,
// pas la peine de revoir le formulaire, on rentre à la maison.
// Retourne true si une redirection a bien été émise.
func RedirectIfAuthed(w http.ResponseWriter, r *http.Request) bool {
	if utils.GetUserIDFromSession(r, db) != 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return true
	}
	return false
}

// currentUser charge l'utilisateur connecté depuis la session, ou renvoie nil.
func currentUser(r *http.Request) *models.User {
	userID := utils.GetUserIDFromSession(r, db)
	if userID == 0 {
		return nil
	}
	u := &models.User{}
	err := db.QueryRow(
		`SELECT id, nom, pseudo, email, photo_url, created_at FROM users WHERE id = ?`,
		userID,
	).Scan(&u.ID, &u.Nom, &u.Pseudo, &u.Email, &u.PhotoURL, &u.CreatedAt)
	if err != nil {
		return nil
	}
	return u
}
