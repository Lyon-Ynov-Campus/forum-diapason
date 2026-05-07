package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"forum-diapason/database"
	"forum-diapason/handlers"
	"forum-diapason/services"
)

// page charge une page et tous les composants réutilisables
func page(name string) *template.Template {
	t := template.Must(template.ParseGlob("./frontend/components/*.html"))
	return template.Must(t.ParseFiles("./frontend/pages/"+name+".html"))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	port   := getEnv("PORT", "8080")
	dbFile := getEnv("DB_FILE", "./forum.db")

	db := database.Init(dbFile)
	defer db.Close()

	handlers.Init(db)

	// Nettoyage automatique des sessions expirées toutes les heures
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			services.DeleteExpiredSessions(db)
		}
	}()

	// ── Pages statiques ──────────────────────────────────────
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
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		page("profile").ExecuteTemplate(w, "profile.html", nil)
	})
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		page("post").ExecuteTemplate(w, "post.html", nil)
	})
	http.Handle("/css/", http.FileServer(http.Dir("./frontend/")))
	http.Handle("/js/",  http.FileServer(http.Dir("./frontend/")))
	http.Handle("/data/", http.FileServer(http.Dir("./frontend/")))

	// ── Auth ─────────────────────────────────────────────────
	http.HandleFunc("/api/auth/register", handlers.Register)
	http.HandleFunc("/api/auth/login",    handlers.Login)
	http.HandleFunc("/api/auth/logout",   handlers.Logout)
	http.HandleFunc("/api/auth/me",       handlers.Me)      // GET=profil, DELETE=supprimer compte

	// ── Posts ────────────────────────────────────────────────
	http.HandleFunc("/api/posts",         handlers.Posts)    // GET=liste, POST=créer
	http.HandleFunc("/api/posts/tag/",    handlers.PostsByTag)
	http.HandleFunc("/api/posts/",        
	func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/like"):
			handlers.LikePost(w, r)
		case strings.HasSuffix(path, "/comments"):
			handlers.Comments(w, r)
		default:
			handlers.PostByID(w, r)
		}
	})

	// ── Comments ─────────────────────────────────────────────
	http.HandleFunc("/api/comments/", handlers.DeleteComment)

	// ── Tags ─────────────────────────────────────────────────
	http.HandleFunc("/api/tags", handlers.Tags)

	// ── Users / Follows ──────────────────────────────────────
	http.HandleFunc("/api/users/", 
	func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/follow"):
			handlers.FollowUser(w, r)
		case strings.HasSuffix(path, "/following"):
			handlers.GetFollowing(w, r)
		case strings.HasSuffix(path, "/followers"):
			handlers.GetFollowers(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	log.Println("Serveur lancé sur http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}