package repositories

import (
	"database/sql"
	"forum-diapason/internal/entities"
)

type UserRepository interface {
	Create(user *entities.User) error
	GetByID(id int) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	GetByPseudo(pseudo string) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entities.User) error {
	query := `INSERT INTO users (email, pseudo, password, photo_url, is_admin, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, user.Email, user.Pseudo, user.Password, user.PhotoURL, user.IsAdmin, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByID(id int) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, email, pseudo, password, photo_url, is_admin, created_at, updated_at FROM users WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Pseudo, &user.Password, &user.PhotoURL, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, email, pseudo, password, photo_url, is_admin, created_at, updated_at FROM users WHERE email = ?`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Pseudo, &user.Password, &user.PhotoURL, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByPseudo(pseudo string) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, email, pseudo, password, photo_url, is_admin, created_at, updated_at FROM users WHERE pseudo = ?`
	err := r.db.QueryRow(query, pseudo).Scan(&user.ID, &user.Email, &user.Pseudo, &user.Password, &user.PhotoURL, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(user *entities.User) error {
	query := `UPDATE users SET email = ?, pseudo = ?, password = ?, photo_url = ?, is_admin = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Email, user.Pseudo, user.Password, user.PhotoURL, user.IsAdmin, user.UpdatedAt, user.ID)
	return err
}

func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
