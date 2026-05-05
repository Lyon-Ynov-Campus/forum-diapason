package models

import "time"

type Post struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Titre           string    `json:"titre"`
	Contenu         string    `json:"contenu"`
	MediaType       string    `json:"media_type"`
	DatePublication time.Time `json:"date_publication"`

	// Champs calculés
	AuthorPseudo string `json:"author_pseudo,omitempty"`
	AuthorPhoto  string `json:"author_photo,omitempty"`
	LikeCount    int    `json:"like_count,omitempty"`
	CommentCount int    `json:"comment_count,omitempty"`
	LikedByMe    bool   `json:"liked_by_me,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

type Comment struct {
	ID      int       `json:"id"`
	PostsID int       `json:"posts_id"`
	UserID  int       `json:"user_id"`
	Contenu string    `json:"contenu"`
	Date    time.Time `json:"date"`

	// Champ calculé
	AuthorPseudo string `json:"author_pseudo,omitempty"`
	AuthorPhoto  string `json:"author_photo,omitempty"`
}

type Tag struct {
	ID  int    `json:"id"`
	Nom string `json:"nom"`
}

type Follow struct {
	FollowerID int `json:"follower_id"`
	FollowedID int `json:"followed_id"`
}