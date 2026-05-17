package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"forum-diapason/database"
	"forum-diapason/handlers"
	"forum-diapason/services"
	"forum-diapason/utils"
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
	utils.CookieSecure = getEnv("COOKIE_SECURE", "false") == "true"

	db := database.Init(dbFile)
	defer db.Close()

	handlers.Init(db)

	// Sinon la table sessions grossit indefiniment
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
	http.HandleFunc("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderPage(w, r, "forgot-password", nil)
	})
	http.HandleFunc("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderPage(w, r, "reset-password", nil)
	})
	http.HandleFunc("/logout", handlers.LogoutPage)
	http.HandleFunc("/profile", handlers.ProfilePage)
	http.HandleFunc("/profile/edit", handlers.ProfileEditPage)
	http.HandleFunc("/profile/avatar", handlers.ProfileAvatarPage)
	http.HandleFunc("/profile/avatar/delete", handlers.ProfileAvatarDeletePage)
	http.HandleFunc("/profile/password", handlers.ProfilePasswordPage)
	http.HandleFunc("/profile/delete", handlers.ProfileDeletePage)
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderPage(w, r, "post", nil)
	})

	// API
	http.HandleFunc("/api/auth/me", handlers.Me)

	// Fichiers statiques
	http.Handle("/css/", http.FileServer(http.Dir("./frontend/")))
	http.Handle("/js/", http.FileServer(http.Dir("./frontend/")))
	http.Handle("/data/", http.FileServer(http.Dir("./frontend/")))
	http.Handle("/avatars/", http.FileServer(http.Dir("./public/")))
	http.Handle("/posts/", http.FileServer(http.Dir("./public/")))
	http.Handle("/image/", http.FileServer(http.Dir("./public/")))

	log.Println("Serveur lancé sur http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
