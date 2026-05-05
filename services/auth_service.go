package services

import (
	"database/sql"
	"errors"
	"forum-diapason/models"
	"forum-diapason/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Register 

func Register(db *sql.DB, nom, pseudo, email, password string) (*models.User, error) {
	// Validations
	if !utils.IsValidEmail(email) {
		return nil, errors.New("email invalide")
	}
	if !utils.IsValidPassword(password) {
		return nil, errors.New("mot de passe trop court (min 8 caractères)")
	}
	if !utils.IsValidPseudo(pseudo) {
		return nil, errors.New("pseudo invalide (3-30 caractères)")
	}

	// Vérifier unicité email + pseudo
	var exists int
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, email).Scan(&exists)
	if exists > 0 {
		return nil, errors.New("email déjà utilisé")
	}
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE pseudo = ?`, pseudo).Scan(&exists)
	if exists > 0 {
		return nil, errors.New("pseudo déjà utilisé")
	}

	// Hasher le mot de passe
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("erreur lors du hashage")
	}

	// Insérer l'utilisateur
	result, err := db.Exec(
		`INSERT INTO users (nom, pseudo, email, password) VALUES (?, ?, ?, ?)`,
		nom, pseudo, email, string(hash),
	)
	if err != nil {
		return nil, errors.New("erreur lors de la création du compte")
	}

	id, _ := result.LastInsertId()
	return &models.User{
		ID:        int(id),
		Nom:       nom,
		Pseudo:    pseudo,
		Email:     email,
		CreatedAt: time.Now(),
	}, nil
}

// Login 

func Login(db *sql.DB, emailOrPseudo, password string) (*models.User, error) {
	user := &models.User{}

	var query string
	if utils.IsValidEmail(emailOrPseudo) {
		query = `SELECT id, nom, pseudo, email, password, photo_url, created_at FROM users WHERE email = ?`
	} else {
		query = `SELECT id, nom, pseudo, email, password, photo_url, created_at FROM users WHERE pseudo = ?`
	}

	err := db.QueryRow(query, emailOrPseudo).Scan(
		&user.ID, &user.Nom, &user.Pseudo, &user.Email,
		&user.Password, &user.PhotoURL, &user.CreatedAt,
	)
	if err != nil {
		return nil, errors.New("identifiants incorrects")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("identifiants incorrects")
	}

	user.Password = "" // ne jamais renvoyer le hash
	return user, nil
}

// CreateSession 

func CreateSession(db *sql.DB, w http.ResponseWriter, userID int) error {
	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return err
	}

	expires := time.Now().Add(utils.SessionDuration)

	_, err = db.Exec(
		`INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`,
		sessionID, userID, expires,
	)
	if err != nil {
		return err
	}

	utils.SetSessionCookie(w, sessionID, expires)
	return nil
}

// Logout 

func Logout(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(utils.SessionCookieName)
	if err == nil {
		db.Exec(`DELETE FROM sessions WHERE id = ?`, cookie.Value)
	}
	utils.ClearSessionCookie(w)
}

// DeleteExpiredSessions 

func DeleteExpiredSessions(db *sql.DB) {
	db.Exec(`DELETE FROM sessions WHERE expires_at <= datetime('now')`)
}
// À appeler périodiquement (ex: goroutine au démarrage)

func DeleteExpiredSessions(db *sql.DB) {
	db.Exec(`DELETE FROM sessions WHERE expires_at <= datetime('now')`)
}