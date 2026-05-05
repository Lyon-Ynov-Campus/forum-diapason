package database

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
	queries := []string{

		`PRAGMA foreign_keys = ON`,

		//USERS
		`CREATE TABLE IF NOT EXISTS users (
			id         INTEGER  PRIMARY KEY AUTOINCREMENT,
			nom        TEXT     NOT NULL,
			pseudo     TEXT     NOT NULL UNIQUE,
			email      TEXT     NOT NULL UNIQUE,
			password   TEXT     NOT NULL,
			photo_url  TEXT     NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		// SESSIONS 
		`CREATE TABLE IF NOT EXISTS sessions (
			id         TEXT     PRIMARY KEY,
			user_id    INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			expires_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user_id    ON sessions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at)`,

		//POSTS
		`CREATE TABLE IF NOT EXISTS posts (
			id               INTEGER  PRIMARY KEY AUTOINCREMENT,
			user_id          INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			titre            TEXT     NOT NULL,
			contenu          TEXT     NOT NULL,
			media_type       TEXT     NOT NULL DEFAULT '',
			date_publication DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_user_id          ON posts(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_posts_date_publication ON posts(date_publication DESC)`,

		// COMMENTS 
		`CREATE TABLE IF NOT EXISTS comments (
			id       INTEGER  PRIMARY KEY AUTOINCREMENT,
			posts_id INTEGER  NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			user_id  INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			contenu  TEXT     NOT NULL,
			date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_posts_id ON comments(posts_id)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_user_id  ON comments(user_id)`,

		// LIKES
		`CREATE TABLE IF NOT EXISTS likes (
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			PRIMARY KEY (user_id, post_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id)`,

		//TAGS 
		`CREATE TABLE IF NOT EXISTS tags (
			id  INTEGER PRIMARY KEY AUTOINCREMENT,
			nom TEXT    NOT NULL UNIQUE COLLATE NOCASE
		)`,

		// POST_TAGS 
		`CREATE TABLE IF NOT EXISTS post_tags (
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			tag_id  INTEGER NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
			PRIMARY KEY (post_id, tag_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_post_tags_tag_id ON post_tags(tag_id)`,

		//FOLLOWS 
		`CREATE TABLE IF NOT EXISTS follows (
			follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			followed_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			PRIMARY KEY (follower_id, followed_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_follows_followed_id ON follows(followed_id)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("Migration échouée: %v", err)
		}
	}
}