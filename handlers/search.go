package handlers

import (
	"database/sql"
	"forum-diapason/models"
	"forum-diapason/services"
	"forum-diapason/utils"
	"net/http"
	"strconv"
	"strings"
)

type SearchResult struct {
	Type   string      `json:"type"` // "user" ou "post"
	User   *models.User `json:"user,omitempty"`
	Post   *models.Post `json:"post,omitempty"`
}

// GET /api/search?q=...
func Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	currentUserID := utils.GetUserIDFromSession(r, db)
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	
	if q == "" {
		sendJSON(w, http.StatusOK, []SearchResult{})
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	results := []SearchResult{}

	// Rechercher les utilisateurs par pseudo
	users, err := searchUsers(db, q, limit)
	if err == nil {
		for _, user := range users {
			results = append(results, SearchResult{
				Type: "user",
				User: user,
			})
		}
	}

	// Rechercher les posts
	posts, err := services.SearchPosts(db, currentUserID, q, "", []string{}, limit, 0)
	if err == nil {
		for _, post := range posts {
			results = append(results, SearchResult{
				Type: "post",
				Post: post,
			})
		}
	}

	sendJSON(w, http.StatusOK, results)
}

// searchUsers cherche les utilisateurs par pseudo
func searchUsers(db *sql.DB, query string, limit int) ([]*models.User, error) {
	like := "%" + strings.ToLower(strings.TrimSpace(query)) + "%"
	rows, err := db.Query(`
		SELECT id, nom, pseudo, email, photo_url, created_at
		FROM users
		WHERE LOWER(pseudo) LIKE ?
		LIMIT ?`,
		like, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		rows.Scan(&u.ID, &u.Nom, &u.Pseudo, &u.Email, &u.PhotoURL, &u.CreatedAt)
		users = append(users, u)
	}
	return users, nil
}
