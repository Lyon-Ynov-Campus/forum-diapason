package entities

import "time"

// Post represents a forum post
type Post struct {
	ID           int       `json:"id"`
	AuthorID     int       `json:"author_id"`
	AuthorPseudo string    `json:"author_pseudo"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Comment represents a post comment
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	AuthorID  int       `json:"author_id"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Like represents a post like
type Like struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
