package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"forum-diapason/database"
	"forum-diapason/handlers"
)

// page charge une page et tous les composants réutilisables
func page(name string) *template.Template {
	t := template.Must(template.ParseGlob("./frontend/components/*.html"))
	return template.Must(t.ParseFiles("./frontend/pages/" + name + ".html"))
}

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
		page("home").ExecuteTemplate(w, "home.html", nil)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		page("login").ExecuteTemplate(w, "login.html", nil)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		page("register").ExecuteTemplate(w, "register.html", nil)
	})
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		page("post").ExecuteTemplate(w, "post.html", nil)
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
