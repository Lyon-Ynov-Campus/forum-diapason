package handlers

// CRUD des posts, likes, commentaires, tags, follows

import (
	"encoding/json"
	"forum-diapason/services"
	"forum-diapason/utils"
	"net/http"
	"strconv"
	"strings"
)

// POSTS

// GET  /api/posts   → liste des posts (public)
// POST /api/posts   → créer un post  (connecté)
func Posts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		currentUserID := utils.GetUserIDFromSession(r, db)

		limit, offset := 20, 0
		if l := r.URL.Query().Get("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
				limit = v
			}
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			if v, err := strconv.Atoi(o); err == nil && v >= 0 {
				offset = v
			}
		}

		posts, err := services.GetPosts(db, currentUserID, limit, offset)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		sendJSON(w, http.StatusOK, posts)

	case http.MethodPost:
		userID := utils.RequireAuth(w, r, db)
		if userID == 0 {
			return
		}
		var body struct {
			Titre     string   `json:"titre"`
			Contenu   string   `json:"contenu"`
			MediaType string   `json:"media_type"`
			Tags      []string `json:"tags"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendError(w, http.StatusBadRequest, "JSON invalide")
			return
		}

		post, err := services.CreatePost(db, userID, body.Titre, body.Contenu, body.MediaType)
		if err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}
		for _, tagNom := range body.Tags {
			if tagNom = strings.TrimSpace(tagNom); tagNom == "" {
				continue
			}
			if tagID, err := services.GetOrCreateTag(db, tagNom); err == nil {
				services.AddTagToPost(db, post.ID, tagID)
			}
		}
		post.Tags = services.GetPostTags(db, post.ID)
		sendJSON(w, http.StatusCreated, post)

	default:
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
	}
}

// GET    /api/posts/{id}  → lire un post
// PUT    /api/posts/{id}  → modifier
// DELETE /api/posts/{id}  → supprimer
func PostByID(w http.ResponseWriter, r *http.Request) {
	postID, err := parseID(r.URL.Path, "/api/posts/")
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID invalide")
		return
	}
	currentUserID := utils.GetUserIDFromSession(r, db)

	switch r.Method {
	case http.MethodGet:
		post, err := services.GetPost(db, postID, currentUserID)
		if err != nil {
			sendError(w, http.StatusNotFound, err.Error())
			return
		}
		sendJSON(w, http.StatusOK, post)

	case http.MethodPut:
		userID := utils.RequireAuth(w, r, db)
		if userID == 0 {
			return
		}
		var body struct {
			Titre   string `json:"titre"`
			Contenu string `json:"contenu"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if err := services.UpdatePost(db, postID, userID, body.Titre, body.Contenu); err != nil {
			code := http.StatusBadRequest
			if err.Error() == "non autorisé" {
				code = http.StatusForbidden
			}
			sendError(w, code, err.Error())
			return
		}
		sendJSON(w, http.StatusOK, map[string]string{"message": "post modifié"})

	case http.MethodDelete:
		userID := utils.RequireAuth(w, r, db)
		if userID == 0 {
			return
		}
		if err := services.DeletePost(db, postID, userID); err != nil {
			code := http.StatusBadRequest
			if err.Error() == "non autorisé" {
				code = http.StatusForbidden
			}
			sendError(w, code, err.Error())
			return
		}
		sendJSON(w, http.StatusOK, map[string]string{"message": "post supprimé"})

	default:
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
	}
}

// LIKES

// POST   /api/posts/{id}/like  → liker
// DELETE /api/posts/{id}/like  → unliker
func LikePost(w http.ResponseWriter, r *http.Request) {
	userID := utils.RequireAuth(w, r, db)
	if userID == 0 {
		return
	}
	postID, err := parseID(r.URL.Path, "/api/posts/")
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID invalide")
		return
	}
	switch r.Method {
	case http.MethodPost:
		if err := services.LikePost(db, userID, postID); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}
		sendJSON(w, http.StatusOK, map[string]string{"message": "liké"})
	case http.MethodDelete:
		if err := services.UnlikePost(db, userID, postID); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}
		sendJSON(w, http.StatusOK, map[string]string{"message": "like retiré"})
	default:
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
	}
}

// COMMENTS

// GET  /api/posts/{id}/comments  → liste des commentaires
// POST /api/posts/{id}/comments  → ajouter
func Comments(w http.ResponseWriter, r *http.Request) {
	postID, err := parseID(r.URL.Path, "/api/posts/")
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID invalide")
		return
	}
	switch r.Method {
	case http.MethodGet:
		comments, err := services.GetComments(db, postID)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		sendJSON(w, http.StatusOK, comments)

	case http.MethodPost:
		userID := utils.RequireAuth(w, r, db)
		if userID == 0 {
			return
		}
		var body struct {
			Contenu string `json:"contenu"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		comment, err := services.CreateComment(db, userID, postID, body.Contenu)
		if err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}
		sendJSON(w, http.StatusCreated, comment)

	default:
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
	}
}

// DELETE /api/comments/{id}
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	userID := utils.RequireAuth(w, r, db)
	if userID == 0 {
		return
	}
	commentID, err := parseID(r.URL.Path, "/api/comments/")
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID invalide")
		return
	}
	if err := services.DeleteComment(db, commentID, userID); err != nil {
		code := http.StatusBadRequest
		if err.Error() == "non autorisé" {
			code = http.StatusForbidden
		}
		sendError(w, code, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, map[string]string{"message": "commentaire supprimé"})
}

// TAGS

// GET /api/tags
func Tags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	tags, err := services.GetAllTags(db)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	sendJSON(w, http.StatusOK, tags)
}

// GET /api/posts/tag/{nom}
func PostsByTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	tagNom := strings.TrimPrefix(r.URL.Path, "/api/posts/tag/")
	if tagNom == "" {
		sendError(w, http.StatusBadRequest, "tag manquant")
		return
	}
	posts, err := services.GetPostsByTag(db, tagNom, utils.GetUserIDFromSession(r, db))
	if err != nil {
		sendError(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	sendJSON(w, http.StatusOK, posts)
}

// FOLLOWS

// POST   /api/users/{id}/follow  → suivre
// DELETE /api/users/{id}/follow  → ne plus suivre
func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID := utils.RequireAuth(w, r, db)
	if followerID == 0 {
		return
	}
	followedID, err := parseID(r.URL.Path, "/api/users/")
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID invalide")
		return
	}
	switch r.Method {
	case http.MethodPost:
		if err := services.Follow(db, followerID, followedID); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}
		sendJSON(w, http.StatusOK, map[string]string{"message": "abonné"})
	case http.MethodDelete:
		if err := services.Unfollow(db, followerID, followedID); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}
		sendJSON(w, http.StatusOK, map[string]string{"message": "désabonné"})
	default:
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
	}
}

// GET /api/users/{id}/following
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	userID, err := parseID(r.URL.Path, "/api/users/")
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID invalide")
		return
	}
	users, err := services.GetFollowing(db, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	sendJSON(w, http.StatusOK, users)
}

// GET /api/users/{id}/followers
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}
	userID, err := parseID(r.URL.Path, "/api/users/")
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID invalide")
		return
	}
	users, err := services.GetFollowers(db, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "erreur serveur")
		return
	}
	sendJSON(w, http.StatusOK, users)
}

// HELPER

func parseID(path, prefix string) (int, error) {
	s := strings.TrimPrefix(path, prefix)
	if idx := strings.Index(s, "/"); idx != -1 {
		s = s[:idx]
	}
	return strconv.Atoi(s)