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

// Register cree un compte
// Les entrees sont nettoyees (espaces, casse) puis validees avant d'ecrire en base
func Register(db *sql.DB, nom, pseudo, email, password string) (*models.User, error) {
	// L'email est normalise en lowercase pour qu'"Alice@Gmail.com" et
	// "alice@gmail.com" pointent sur le meme compte
	// Le pseudo lui reste sensible a la casse (choix de design : "Alice" ≠ "alice")
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

	// On verifie l'unicite explicitement pour pouvoir donner un message cible
	// ("email deja utilise" plutot qu'un echec d'INSERT generique)
	// La contrainte UNIQUE de la DB sert quand meme de filet en cas de course concurrente
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

// Login accepte indifferemment un email ou un pseudo dans le meme champ
// Le message d'erreur reste identique (compte introuvable / mauvais mdp)
// pour ne pas reveler quels comptes existent
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

// CreateSession genere un token, l'enregistre en base et pose le cookie cote client
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

// Logout efface la session cote DB et cote cookie
// Si le cookie est absent (user deja deco), on ne fait rien de bruyant
func Logout(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(utils.SessionCookieName)
	if err == nil {
		db.Exec(`DELETE FROM sessions WHERE id = ?`, cookie.Value)
	}
	utils.ClearSessionCookie(w)
}

// DeleteExpiredSessions vide la table des sessions perimees en une requete
func DeleteExpiredSessions(db *sql.DB) {
	db.Exec(`DELETE FROM sessions WHERE expires_at <= datetime('now')`)
}

// StartSessionCleanup lance une goroutine qui purge les sessions expirees
// a intervalles reguliers, pour la duree de vie du programme
// A appeler une fois au demarrage
func StartSessionCleanup(db *sql.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			DeleteExpiredSessions(db)
		}
	}()
}

func CreatePasswordReset(db *sql.DB, email string) (string, int, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	var userID int
	if err := db.QueryRow(`SELECT id FROM users WHERE email = ?`, email).Scan(&userID); err != nil {
		return "", 0, errors.New("email inconnu")
	}
	token, err := utils.GenerateSessionID()
	if err != nil {
		return "", 0, err
	}
	expires := time.Now().Add(1 * time.Hour)
	_, err = db.Exec(
		`INSERT INTO password_resets (token, user_id, expires_at) VALUES (?, ?, ?)`,
		token, userID, expires,
	)
	if err != nil {
		return "", 0, err
	}
	return token, userID, nil
}

func ResetPasswordWithToken(db *sql.DB, token, newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("mot de passe trop court (8 caractères min)")
	}
	var userID int
	var used int
	err := db.QueryRow(
		`SELECT user_id, used FROM password_resets WHERE token = ? AND expires_at > datetime('now')`,
		token,
	).Scan(&userID, &used)
	if err != nil {
		return errors.New("token invalide ou expiré")
	}
	if used == 1 {
		return errors.New("token déjà utilisé")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE users SET password = ? WHERE id = ?`, string(hash), userID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`UPDATE password_resets SET used = 1 WHERE token = ?`, token); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
