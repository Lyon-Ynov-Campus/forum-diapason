package database

// Ajouter les tables dans runMigrations() au fur et à mesure.

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Init(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Impossible d'ouvrir la base de données: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Impossible de contacter la base de données: %v", err)
	}
	runMigrations(db)
	return db
}

func runMigrations(db *sql.DB) {
	// Ajouter les CREATE TABLE IF NOT EXISTS ici
	queries := []string{}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("Migration échouée: %v", err)
		}
	}
}
