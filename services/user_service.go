package services

// Modifs du profil user : infos, avatar, mdp

import (
	"database/sql"
	"errors"
	"fmt"
	"image"
	_ "image/gif"  // decoders enregistres via side-effect import
	_ "image/jpeg" //
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"forum-diapason/utils"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // decoder WEBP (stdlib n'en a pas)
)

// UpdateProfile met a jour nom/pseudo/email
// Les checks d'unicite excluent l'user lui-meme (sinon il pourrait pas garder
// son propre pseudo en editant juste son nom)
// La photo de profil est geree a part via UpdateAvatar
func UpdateProfile(db *sql.DB, userID int, nom, pseudo, email string) error {
	nom = strings.TrimSpace(nom)
	pseudo = strings.TrimSpace(pseudo)
	email = strings.ToLower(strings.TrimSpace(email))

	if !utils.IsValidNom(nom) {
		return errors.New("nom invalide (1-50 caractères)")
	}
	if !utils.IsValidPseudo(pseudo) {
		return errors.New("pseudo invalide (3-30 caractères, lettres/chiffres/-/_)")
	}
	if !utils.IsValidEmail(email) {
		return errors.New("email invalide")
	}

	var exists int
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ? AND id != ?`, email, userID).Scan(&exists)
	if exists > 0 {
		return errors.New("email déjà utilisé")
	}
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE pseudo = ? AND id != ?`, pseudo, userID).Scan(&exists)
	if exists > 0 {
		return errors.New("pseudo déjà utilisé")
	}

	_, err := db.Exec(
		`UPDATE users SET nom = ?, pseudo = ?, email = ? WHERE id = ?`,
		nom, pseudo, email, userID,
	)
	if err != nil {
		return errors.New("erreur lors de la mise à jour du profil")
	}
	return nil
}

// UpdatePassword change le mdp apres avoir verifie l'ancien
// Le hash est regenere avec bcrypt comme a l'inscription
func UpdatePassword(db *sql.DB, userID int, oldPassword, newPassword string) error {
	if !utils.IsValidPassword(newPassword) {
		return errors.New("nouveau mot de passe trop court (min 8 caractères)")
	}

	var currentHash string
	err := db.QueryRow(`SELECT password FROM users WHERE id = ?`, userID).Scan(&currentHash)
	if err != nil {
		return errors.New("utilisateur introuvable")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(oldPassword)); err != nil {
		return errors.New("ancien mot de passe incorrect")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("erreur lors du hashage")
	}

	_, err = db.Exec(`UPDATE users SET password = ? WHERE id = ?`, string(newHash), userID)
	if err != nil {
		return errors.New("erreur lors de la mise à jour du mot de passe")
	}
	return nil
}

// --- Avatar ---

const (
	AvatarMaxSize = 2 << 20 // 2 MiB en upload
	AvatarSize    = 256     // px, sortie carree
	AvatarDir     = "./public/avatars"
)

// allowedAvatarMIME : on accepte ces formats en entree
// La sortie est toujours du PNG (transparence preservee, format unique sur disque)
var allowedAvatarMIME = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

// UpdateAvatar recoit un fichier image, le valide, le resize en carre
// AvatarSize × AvatarSize, l'enregistre en PNG sous public/avatars/{userID}.png
// et met a jour photo_url en base
func UpdateAvatar(db *sql.DB, userID int, file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Size > AvatarMaxSize {
		return "", errors.New("image trop lourde (max 2 Mo)")
	}

	// Detection MIME via les 512 premiers octets (l'extension du nom de fichier
	// peut mentir, un user peut renommer un .exe en .png)
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	mime := http.DetectContentType(buf[:n])
	if !allowedAvatarMIME[mime] {
		return "", errors.New("format non supporté (JPG, PNG, WEBP, GIF)")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", errors.New("erreur de lecture du fichier")
	}

	// Decode generique : image.Decode utilise les decoders enregistres en
	// haut du fichier (jpeg/png/gif/webp)
	src, _, err := image.Decode(file)
	if err != nil {
		return "", errors.New("image illisible")
	}

	resized := resizeToSquare(src, AvatarSize)

	if err := os.MkdirAll(AvatarDir, 0755); err != nil {
		return "", errors.New("erreur lors de la création du dossier d'avatars")
	}

	// On sort toujours en PNG
	// Si l'user avait un fichier dans un autre format (legacy), on le supprime
	// pour pas laisser d'orphelin
	for _, oldExt := range []string{".jpg", ".jpeg", ".webp", ".gif"} {
		os.Remove(filepath.Join(AvatarDir, fmt.Sprintf("%d%s", userID, oldExt)))
	}

	filename := fmt.Sprintf("%d.png", userID)
	fullPath := filepath.Join(AvatarDir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", errors.New("erreur lors de l'enregistrement")
	}
	defer dst.Close()

	if err := png.Encode(dst, resized); err != nil {
		return "", errors.New("erreur lors de l'écriture")
	}

	// Le ?v=timestamp force le browser a recharger apres update (sinon il
	// garde la version en cache, le chemin n'ayant pas change)
	webPath := fmt.Sprintf("/avatars/%s?v=%d", filename, time.Now().Unix())
	if _, err := db.Exec(`UPDATE users SET photo_url = ? WHERE id = ?`, webPath, userID); err != nil {
		return "", errors.New("erreur lors de la mise à jour du profil")
	}
	return webPath, nil
}

// DeleteAvatar supprime le fichier sur disque et vide photo_url en base
// Idempotent : si l'user n'a pas d'avatar, on ne renvoie pas d'erreur
func DeleteAvatar(db *sql.DB, userID int) error {
	// On tente toutes les extensions au cas ou des fichiers legacy traineraient
	for _, ext := range []string{".png", ".jpg", ".jpeg", ".webp", ".gif"} {
		os.Remove(filepath.Join(AvatarDir, fmt.Sprintf("%d%s", userID, ext)))
	}
	if _, err := db.Exec(`UPDATE users SET photo_url = '' WHERE id = ?`, userID); err != nil {
		return errors.New("erreur lors de la mise à jour du profil")
	}
	return nil
}

// resizeToSquare prend une image, la center-crop en carre et la resize
// avec un filtre de bonne qualite (CatmullRom)
// Garantit une sortie size × size quelle que soit la resolution/le ratio d'entree
func resizeToSquare(src image.Image, size int) image.Image {
	b := src.Bounds()
	side := b.Dx()
	if b.Dy() < side {
		side = b.Dy()
	}
	// rectangle carre centre dans la source
	offX := (b.Dx() - side) / 2
	offY := (b.Dy() - side) / 2
	sr := image.Rect(b.Min.X+offX, b.Min.Y+offY, b.Min.X+offX+side, b.Min.Y+offY+side)

	out := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(out, out.Bounds(), src, sr, draw.Over, nil)
	return out
}
