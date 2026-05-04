package repositories

import (
	"database/sql"
	"forum-diapason/internal/domain/entities"
)

type PostRepository interface {
	Create(post *entities.Post) error
	GetByID(id int) (*entities.Post, error)
	GetAll(limit, offset int) ([]*entities.Post, error)
	GetByAuthorID(authorID int, limit, offset int) ([]*entities.Post, error)
	Update(post *entities.Post) error
	Delete(id int) error
	Like(postID, userID int) error
	Unlike(postID, userID int) error
	GetLikeCount(postID int) (int, error)
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *entities.Post) error {
	query := `INSERT INTO posts (author_id, author_pseudo, title, content, like_count, comment_count, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, post.AuthorID, post.AuthorPseudo, post.Title, post.Content, post.LikeCount, post.CommentCount, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = int(id)
	return nil
}

func (r *postRepository) GetByID(id int) (*entities.Post, error) {
	post := &entities.Post{}
	query := `SELECT id, author_id, author_pseudo, title, content, like_count, comment_count, created_at, updated_at FROM posts WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&post.ID, &post.AuthorID, &post.AuthorPseudo, &post.Title, &post.Content, &post.LikeCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *postRepository) GetAll(limit, offset int) ([]*entities.Post, error) {
	query := `SELECT id, author_id, author_pseudo, title, content, like_count, comment_count, created_at, updated_at FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		post := &entities.Post{}
		err := rows.Scan(&post.ID, &post.AuthorID, &post.AuthorPseudo, &post.Title, &post.Content, &post.LikeCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *postRepository) GetByAuthorID(authorID int, limit, offset int) ([]*entities.Post, error) {
	query := `SELECT id, author_id, author_pseudo, title, content, like_count, comment_count, created_at, updated_at FROM posts WHERE author_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		post := &entities.Post{}
		err := rows.Scan(&post.ID, &post.AuthorID, &post.AuthorPseudo, &post.Title, &post.Content, &post.LikeCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *postRepository) Update(post *entities.Post) error {
	query := `UPDATE posts SET title = ?, content = ?, like_count = ?, comment_count = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, post.Title, post.Content, post.LikeCount, post.CommentCount, post.UpdatedAt, post.ID)
	return err
}

func (r *postRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *postRepository) Like(postID, userID int) error {
	query := `INSERT OR IGNORE INTO likes (post_id, user_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, postID, userID)
	if err != nil {
		return err
	}
	// Update like count
	return r.updateLikeCount(postID)
}

func (r *postRepository) Unlike(postID, userID int) error {
	query := `DELETE FROM likes WHERE post_id = ? AND user_id = ?`
	_, err := r.db.Exec(query, postID, userID)
	if err != nil {
		return err
	}
	// Update like count
	return r.updateLikeCount(postID)
}

func (r *postRepository) updateLikeCount(postID int) error {
	query := `UPDATE posts SET like_count = (SELECT COUNT(*) FROM likes WHERE post_id = ?) WHERE id = ?`
	_, err := r.db.Exec(query, postID, postID)
	return err
}

func (r *postRepository) GetLikeCount(postID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM likes WHERE post_id = ?`
	err := r.db.QueryRow(query, postID).Scan(&count)
	return count, err
}
