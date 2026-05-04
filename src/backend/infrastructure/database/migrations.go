package database

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			pseudo TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			photo_url TEXT,
			is_admin BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author_id INTEGER NOT NULL,
			author_pseudo TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			like_count INTEGER DEFAULT 0,
			comment_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			author_id INTEGER NOT NULL,
			author TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(post_id, user_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id)`,
		`CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id)`,
		`CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id)`,
	}

	for i, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("failed to run migration %d: %v", i+1, err)
		}
	}

	return nil
}
