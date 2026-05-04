package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// NewSQLiteDB creates and returns a new SQLite database connection
func NewSQLiteDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}
