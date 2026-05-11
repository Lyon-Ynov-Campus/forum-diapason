package main

import (
	"net/http"
	"strings"

	"forum-diapason/handlers"
	"forum-diapason/services"
	"forum-diapason/utils"
)

func postsRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case strings.HasSuffix(path, "/top"):
		handlers.TopPosts(w, r)
	case strings.Contains(path, "/like"):
		handlers.LikePost(w, r)
	case strings.Contains(path, "/comments"):
		handlers.Comments(w, r)
	case strings.Contains(path, "/tag/"):
		handlers.PostsByTag(w, r)
	default:
		handlers.PostByID(w, r)
	}
}

func usersRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case strings.HasSuffix(path, "/follow"):
		handlers.FollowUser(w, r)
	case strings.HasSuffix(path, "/following"):
		handlers.GetFollowing(w, r)
	case strings.HasSuffix(path, "/followers"):
		handlers.GetFollowers(w, r)
	case strings.HasSuffix(path, "/posts"):
		userID := parseID(path, "/api/users/")
		currentUserID := utils.GetUserIDFromSession(r, db)
		posts, err := services.GetPostsByUser(db, userID, currentUserID, 20, 0)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erreur serveur")
			return
		}
		writeJSON(w, http.StatusOK, posts)
	default:
		userID := parseID(path, "/api/users/")
		user, err := services.GetUserByID(db, userID)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, user)
	}
}
