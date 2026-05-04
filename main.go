package main

import (
	"log"
	"net/http"
	"os"

	"forum-diapason/database"
	"forum-diapason/handlers"
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

	// Pages
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./frontend/index.html")
	})

	// API — à compléter
	// http.HandleFunc("/api/auth/register", handlers.Register)
	// http.HandleFunc("/api/posts", handlers.Posts)

	// Fichiers statiques
	http.Handle("/css/", http.FileServer(http.Dir("./frontend/")))
	http.Handle("/js/", http.FileServer(http.Dir("./frontend/")))

	log.Println("Serveur lancé sur http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
