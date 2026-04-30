package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	frontendDir, err := resolveFrontendDir()
	if err != nil {
		log.Fatalf("failed to resolve frontend directory: %v", err)
	}

	fileServer := http.FileServer(http.Dir(frontendDir))
	mux := http.NewServeMux()

	// Serve static assets directly.
	mux.Handle("/css/", fileServer)
	mux.Handle("/js/", fileServer)

	// Serve the application entry point for the root route.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
	})

	addr := ":3000"
	log.Printf("Frontend server running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("frontend server stopped: %v", err)
	}
}

func resolveFrontendDir() (string, error) {
	candidates := []string{
		"./frontend",
		"../frontend",
	}

	for _, candidate := range candidates {
		indexPath := filepath.Join(candidate, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			return candidate, nil
		}
	}

	return "", os.ErrNotExist
}
