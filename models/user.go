package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Nom       string    `json:"nom"`
	Pseudo    string    `json:"pseudo"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	PhotoURL  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}