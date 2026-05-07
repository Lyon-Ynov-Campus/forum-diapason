package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"forum-diapason/database"
	"forum-diapason/handlers"
	"forum-diapason/services"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	port := getEnv("PORT", "8080")
	dbFile := getEnv("DB_FILE", "./forum.db")

	db := database.Init(dbFile)
	defer db.Close()

	handlers.Init(db)

	// Sinon la table sessions grossit indéfiniment
	services.StartSessionCleanup(db, time.Hour)

	// Pages
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		handlers.RenderPage(w, r, "home", nil)
	})
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/register", handlers.RegisterPage)
	http.HandleFunc("/logout", handlers.LogoutPage)
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		if handlers.RequirePageAuth(w, r) == 0 {
			return
		}
		handlers.RenderPage(w, r, "profile", nil)
	})
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderPage(w, r, "post", nil)
	})

	// API
	http.HandleFunc("/api/auth/me", handlers.Me)

	// Fichiers statiques
	http.Handle("/css/", http.FileServer(http.Dir("./frontend/")))
	http.Handle("/js/", http.FileServer(http.Dir("./frontend/")))
	http.Handle("/data/", http.FileServer(http.Dir("./frontend/")))

	log.Println("Serveur lancé sur http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
