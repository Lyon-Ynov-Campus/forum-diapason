package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"forum-diapason/database"
	"forum-diapason/handlers"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
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

	db := database.Init(dbFile)
	defer db.Close()

	handlers.Init(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/posts", handlers.Posts)
	mux.HandleFunc("/api/posts/", postsRouter)

	mux.HandleFunc("/api/comments/", handlers.DeleteComment)

	mux.HandleFunc("/api/tags", handlers.Tags)

	log.Println("API démarrée sur http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, cors(mux)))
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
