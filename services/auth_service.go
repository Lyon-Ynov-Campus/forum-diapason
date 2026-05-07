package services

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"forum-diapason/models"
	"forum-diapason/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register crée un compte. Les entrées sont nettoyées (espaces, casse) puis
// validées avant d'écrire en base.
func Register(db *sql.DB, nom, pseudo, email, password string) (*models.User, error) {
	// L'email est normalisé en minuscules pour qu'"Alice@Gmail.com" et
	// "alice@gmail.com" pointent sur le même compte. Le pseudo, lui, reste
	// sensible à la casse (choix de design : "Alice" ≠ "alice").
	nom = strings.TrimSpace(nom)
	pseudo = strings.TrimSpace(pseudo)
	email = strings.ToLower(strings.TrimSpace(email))

	if !utils.IsValidNom(nom) {
		return nil, errors.New("nom invalide (1-50 caractères)")
	}
	if !utils.IsValidPseudo(pseudo) {
		return nil, errors.New("pseudo invalide (3-30 caractères, lettres/chiffres/-/_)")
	}
	if !utils.IsValidEmail(email) {
		return nil, errors.New("email invalide")
	}
	if !utils.IsValidPassword(password) {
		return nil, errors.New("mot de passe trop court (min 8 caractères)")
	}

	// On vérifie l'unicité explicitement pour pouvoir donner un message ciblé
	// ("email déjà utilisé" plutôt qu'un échec d'INSERT générique). La contrainte
	// UNIQUE de la DB sert quand même de filet en cas de course concurrente.
	var exists int
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, email).Scan(&exists)
	if exists > 0 {
		return nil, errors.New("email déjà utilisé")
	}
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE pseudo = ?`, pseudo).Scan(&exists)
	if exists > 0 {
		return nil, errors.New("pseudo déjà utilisé")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("erreur lors du hashage")
	}

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

// Login accepte indifféremment un email ou un pseudo dans le même champ.
// Le message d'erreur reste identique (compte introuvable / mauvais mot de passe)
// pour ne pas révéler quels comptes existent.
func Login(db *sql.DB, emailOrPseudo, password string) (*models.User, error) {
	emailOrPseudo = strings.TrimSpace(emailOrPseudo)
	user := &models.User{}

	var query string
	if utils.IsValidEmail(emailOrPseudo) {
		emailOrPseudo = strings.ToLower(emailOrPseudo)
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

	user.Password = "" // on ne fait jamais sortir le hash de cette fonction
	return user, nil
}

// CreateSession génère un token, l'enregistre en base et pose le cookie côté client.
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

// Logout efface la session côté DB et côté cookie. Si le cookie est absent
// (utilisateur déjà déconnecté), on ne fait rien de bruyant.
func Logout(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(utils.SessionCookieName)
	if err == nil {
		db.Exec(`DELETE FROM sessions WHERE id = ?`, cookie.Value)
	}
	utils.ClearSessionCookie(w)
}

// DeleteExpiredSessions vide la table des sessions périmées en une requête.
func DeleteExpiredSessions(db *sql.DB) {
	db.Exec(`DELETE FROM sessions WHERE expires_at <= datetime('now')`)
}

// StartSessionCleanup lance une goroutine qui purge les sessions expirées à
// intervalles réguliers, pour la durée de vie du programme. À appeler une fois
// au démarrage.
func StartSessionCleanup(db *sql.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			DeleteExpiredSessions(db)
		}
	}()
}
