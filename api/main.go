package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum-diapason/database"
	"forum-diapason/handlers"
)

var db *sql.DB

func main() {
	port   := getEnv("API_PORT", "8081")
	dbFile := getEnv("DB_FILE", "./forum.db")

	db = database.Init(dbFile)
	defer db.Close()

	handlers.Init(db)

	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("/api/auth/register", authRegister)
	mux.HandleFunc("/api/auth/login",    authLogin)
	mux.HandleFunc("/api/auth/logout",   authLogout)
	mux.HandleFunc("/api/auth/me",       authMe)

	// Profil
	mux.HandleFunc("/api/profile",          profileUpdate)
	mux.HandleFunc("/api/profile/password", profilePassword)

	// Posts
	mux.HandleFunc("/api/posts",  handlers.Posts)
	mux.HandleFunc("/api/posts/", postsRouter)

	// Comments
	mux.HandleFunc("/api/comments/", handlers.DeleteComment)

	// Tags
	mux.HandleFunc("/api/tags", handlers.Tags)

	// Users
	mux.HandleFunc("/api/users/", usersRouter)

	log.Println("API démarrée sur http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, cors(mux)))
}
