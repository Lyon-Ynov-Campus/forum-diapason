package entities

import "time"

// User represents a forum user
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Pseudo    string    `json:"pseudo"`
	Password  string    `json:"-"` // Never expose password in JSON
	PhotoURL  string    `json:"photo_url"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new User entity
func NewUser(email, pseudo, password string) *User {
	return &User{
		Email:     email,
		Pseudo:    pseudo,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
