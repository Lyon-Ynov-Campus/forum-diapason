package database

import (
	"database/sql"
	"fmt"

	"forum-diapason/pkg/config"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteConnection(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.Database.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Run migrations
	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return db, nil
}
