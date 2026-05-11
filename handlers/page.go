package handlers

// Helpers partages par tous les handlers de pages : rendu de template
// + gardes d'auth pensees pour le HTML (redirects plutot que JSON)

import (
	"html/template"
	"net/http"

	"forum-diapason/models"
	"forum-diapason/utils"
)

// RenderPage charge la page demandee + les composants partages (header, footer...)
// et l'execute
// On glisse au passage le user connecte sous la cle "User" pour que le header
// sache afficher le pseudo (ou les boutons sign in/register)
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

// RequirePageAuth est la version "page" de utils.RequireAuth
// Si le user n'est pas co, on redirige vers /login plutot que de repondre 401 en JSON
// Retourne 0 si la redirection a ete emise, sinon l'userID
func RequirePageAuth(w http.ResponseWriter, r *http.Request) int {
	userID := utils.GetUserIDFromSession(r, db)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 0
	}
	return userID
}

// RedirectIfAuthed sert sur /login et /register : si on est deja co,
// pas la peine de revoir le form, on rentre a la maison
// Retourne true si une redirection a bien ete emise
func RedirectIfAuthed(w http.ResponseWriter, r *http.Request) bool {
	if utils.GetUserIDFromSession(r, db) != 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return true
	}
	return false
}

// currentUser charge l'user co depuis la session, ou renvoie nil
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
