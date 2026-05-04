package handlers

// handlers/ — Un fichier par feature (auth.go, post.go, etc.)
// Chaque handler reçoit (w http.ResponseWriter, r *http.Request) et répond en JSON.

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

var db *sql.DB

func Init(database *sql.DB) {
	db = database
}

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJSON(w, status, map[string]string{"error": msg})
}
