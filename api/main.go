package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"forum-diapason/database"
	"forum-diapason/handlers"
	"forum-diapason/services"
	"forum-diapason/utils"
)

var db *sql.DB

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", getEnv("FRONTEND_ORIGIN", "http://localhost:8080"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := getEnv("API_PORT", "8081")
	dbFile := getEnv("DB_FILE", "./forum.db")

	db = database.Init(dbFile)
	defer db.Close()

	handlers.Init(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/register", authRegister)
	mux.HandleFunc("/api/auth/login",    authLogin)
	mux.HandleFunc("/api/auth/logout",   authLogout)
	mux.HandleFunc("/api/auth/me",       authMe)

	mux.HandleFunc("/api/posts", handlers.Posts)
	mux.HandleFunc("/api/posts/", postsRouter)

	mux.HandleFunc("/api/comments/", handlers.DeleteComment)

	mux.HandleFunc("/api/tags", handlers.Tags)

	log.Println("API démarrée sur http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, cors(mux)))
}

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

func postsRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
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
