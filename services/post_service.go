package services

import (
	"database/sql"
	"errors"
	"forum-diapason/models"
	"strings"
	"time"
)


// POSTS

func CreatePost(db *sql.DB, userID int, titre, contenu, mediaType string) (*models.Post, error) {
	if strings.TrimSpace(titre) == "" {
		return nil, errors.New("le titre est obligatoire")
	}
	if strings.TrimSpace(contenu) == "" {
		return nil, errors.New("le contenu est obligatoire")
	}

	result, err := db.Exec(
		`INSERT INTO posts (user_id, titre, contenu, media_type) VALUES (?, ?, ?, ?)`,
		userID, strings.TrimSpace(titre), strings.TrimSpace(contenu), mediaType,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.Post{
		ID:              int(id),
		UserID:          userID,
		Titre:           titre,
		Contenu:         contenu,
		MediaType:       mediaType,
		DatePublication: time.Now(),
	}, nil
}

func GetPost(db *sql.DB, postID, currentUserID int) (*models.Post, error) {
	post := &models.Post{}
	query := `
		SELECT p.id, p.user_id, p.titre, p.contenu, p.media_type, p.date_publication,
		       u.pseudo, u.photo_url,
		       (SELECT COUNT(*) FROM likes    WHERE post_id  = p.id) AS like_count,
		       (SELECT COUNT(*) FROM comments WHERE posts_id = p.id) AS comment_count,
		       (SELECT COUNT(*) FROM likes    WHERE post_id  = p.id AND user_id = ?) AS liked_by_me
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.id = ?`

	err := db.QueryRow(query, currentUserID, postID).Scan(
		&post.ID, &post.UserID, &post.Titre, &post.Contenu,
		&post.MediaType, &post.DatePublication,
		&post.AuthorPseudo, &post.AuthorPhoto,
		&post.LikeCount, &post.CommentCount, &post.LikedByMe,
	)
	if err != nil {
		return nil, errors.New("post introuvable")
	}

	post.Tags = GetPostTags(db, postID)
	return post, nil
}

func GetPosts(db *sql.DB, currentUserID, limit, offset int) ([]*models.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.titre, p.contenu, p.media_type, p.date_publication,
		       u.pseudo, u.photo_url,
		       (SELECT COUNT(*) FROM likes    WHERE post_id  = p.id) AS like_count,
		       (SELECT COUNT(*) FROM comments WHERE posts_id = p.id) AS comment_count,
		       (SELECT COUNT(*) FROM likes    WHERE post_id  = p.id AND user_id = ?) AS liked_by_me
		FROM posts p
		JOIN users u ON u.id = p.user_id
		ORDER BY p.date_publication DESC
		LIMIT ? OFFSET ?`

	rows, err := db.Query(query, currentUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Titre, &p.Contenu,
			&p.MediaType, &p.DatePublication,
			&p.AuthorPseudo, &p.AuthorPhoto,
			&p.LikeCount, &p.CommentCount, &p.LikedByMe,
		)
		if err != nil {
			continue
		}
		p.Tags = GetPostTags(db, p.ID)
		posts = append(posts, p)
	}
	return posts, nil
}

func UpdatePost(db *sql.DB, postID, userID int, titre, contenu string) error {
	// Vérifier que le post appartient à l'utilisateur
	var ownerID int
	err := db.QueryRow(`SELECT user_id FROM posts WHERE id = ?`, postID).Scan(&ownerID)
	if err != nil {
		return errors.New("post introuvable")
	}
	if ownerID != userID {
		return errors.New("non autorisé")
	}

	_, err = db.Exec(
		`UPDATE posts SET titre = ?, contenu = ? WHERE id = ?`,
		strings.TrimSpace(titre), strings.TrimSpace(contenu), postID,
	)
	return err
}

func DeletePost(db *sql.DB, postID, userID int) error {
	var ownerID int
	err := db.QueryRow(`SELECT user_id FROM posts WHERE id = ?`, postID).Scan(&ownerID)
	if err != nil {
		return errors.New("post introuvable")
	}
	if ownerID != userID {
		return errors.New("non autorisé")
	}
	// CASCADE supprime automatiquement les comments, likes et post_tags liés
	_, err = db.Exec(`DELETE FROM posts WHERE id = ?`, postID)
	return err
}


// LIKES


func LikePost(db *sql.DB, userID, postID int) error {
	var exists int
	db.QueryRow(`SELECT COUNT(*) FROM posts WHERE id = ?`, postID).Scan(&exists)
	if exists == 0 {
		return errors.New("post introuvable")
	}
	// INSERT OR IGNORE : si déjà liké, ne fait rien (pas d'erreur)
	_, err := db.Exec(
		`INSERT OR IGNORE INTO likes (user_id, post_id) VALUES (?, ?)`,
		userID, postID,
	)
	return err
}

func UnlikePost(db *sql.DB, userID, postID int) error {
	_, err := db.Exec(
		`DELETE FROM likes WHERE user_id = ? AND post_id = ?`,
		userID, postID,
	)
	return err
}

// ═══════════════════════════════════════════════════════════════════
// COMMENTS
// ═══════════════════════════════════════════════════════════════════

func CreateComment(db *sql.DB, userID, postID int, contenu string) (*models.Comment, error) {
	if strings.TrimSpace(contenu) == "" {
		return nil, errors.New("le commentaire est vide")
	}

	// Vérifier que le post existe
	var exists int
	db.QueryRow(`SELECT COUNT(*) FROM posts WHERE id = ?`, postID).Scan(&exists)
	if exists == 0 {
		return nil, errors.New("post introuvable")
	}

	result, err := db.Exec(
		`INSERT INTO comments (posts_id, user_id, contenu) VALUES (?, ?, ?)`,
		postID, userID, strings.TrimSpace(contenu),
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	// Récupérer le pseudo pour la réponse
	var pseudo, photo string
	db.QueryRow(`SELECT pseudo, photo_url FROM users WHERE id = ?`, userID).Scan(&pseudo, &photo)

	return &models.Comment{
		ID:           int(id),
		PostsID:      postID,
		UserID:       userID,
		Contenu:      contenu,
		Date:         time.Now(),
		AuthorPseudo: pseudo,
		AuthorPhoto:  photo,
	}, nil
}

func GetComments(db *sql.DB, postID int) ([]*models.Comment, error) {
	rows, err := db.Query(`
		SELECT c.id, c.posts_id, c.user_id, c.contenu, c.date,
		       u.pseudo, u.photo_url
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.posts_id = ?
		ORDER BY c.date ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		rows.Scan(&c.ID, &c.PostsID, &c.UserID, &c.Contenu, &c.Date,
			&c.AuthorPseudo, &c.AuthorPhoto)
		comments = append(comments, c)
	}
	return comments, nil
}

func DeleteComment(db *sql.DB, commentID, userID int) error {
	var ownerID int
	err := db.QueryRow(`SELECT user_id FROM comments WHERE id = ?`, commentID).Scan(&ownerID)
	if err != nil {
		return errors.New("commentaire introuvable")
	}
	if ownerID != userID {
		return errors.New("non autorisé")
	}
	_, err = db.Exec(`DELETE FROM comments WHERE id = ?`, commentID)
	return err
}


// TAGS


// GetOrCreateTag retourne l'id d'un tag existant ou en crée un nouveau
func GetOrCreateTag(db *sql.DB, nom string) (int, error) {
	nom = strings.TrimSpace(strings.ToLower(nom))
	if nom == "" {
		return 0, errors.New("nom de tag vide")
	}

	var id int
	err := db.QueryRow(`SELECT id FROM tags WHERE nom = ?`, nom).Scan(&id)
	if err == nil {
		return id, nil // tag existe déjà
	}

	result, err := db.Exec(`INSERT INTO tags (nom) VALUES (?)`, nom)
	if err != nil {
		return 0, err
	}
	newID, _ := result.LastInsertId()
	return int(newID), nil
}

func AddTagToPost(db *sql.DB, postID, tagID int) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO post_tags (post_id, tag_id) VALUES (?, ?)`,
		postID, tagID,
	)
	return err
}

func RemoveTagFromPost(db *sql.DB, postID, tagID int) error {
	_, err := db.Exec(
		`DELETE FROM post_tags WHERE post_id = ? AND tag_id = ?`,
		postID, tagID,
	)
	return err
}

func GetPostTags(db *sql.DB, postID int) []string {
	rows, err := db.Query(`
		SELECT t.nom FROM tags t
		JOIN post_tags pt ON pt.tag_id = t.id
		WHERE pt.post_id = ?`, postID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var nom string
		rows.Scan(&nom)
		tags = append(tags, nom)
	}
	return tags
}

func GetAllTags(db *sql.DB) ([]*models.Tag, error) {
	rows, err := db.Query(`SELECT id, nom FROM tags ORDER BY nom`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		t := &models.Tag{}
		rows.Scan(&t.ID, &t.Nom)
		tags = append(tags, t)
	}
	return tags, nil
}

func GetPostsByTag(db *sql.DB, tagNom string, currentUserID int) ([]*models.Post, error) {
	rows, err := db.Query(`
		SELECT p.id, p.user_id, p.titre, p.contenu, p.media_type, p.date_publication,
		       u.pseudo, u.photo_url,
		       (SELECT COUNT(*) FROM likes    WHERE post_id  = p.id) AS like_count,
		       (SELECT COUNT(*) FROM comments WHERE posts_id = p.id) AS comment_count,
		       (SELECT COUNT(*) FROM likes    WHERE post_id  = p.id AND user_id = ?) AS liked_by_me
		FROM posts p
		JOIN users u ON u.id = p.user_id
		JOIN post_tags pt ON pt.post_id = p.id
		JOIN tags t ON t.id = pt.tag_id
		WHERE t.nom = ?
		ORDER BY p.date_publication DESC`, currentUserID, tagNom)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}
		rows.Scan(
			&p.ID, &p.UserID, &p.Titre, &p.Contenu,
			&p.MediaType, &p.DatePublication,
			&p.AuthorPseudo, &p.AuthorPhoto,
			&p.LikeCount, &p.CommentCount, &p.LikedByMe,
		)
		p.Tags = GetPostTags(db, p.ID)
		posts = append(posts, p)
	}
	return posts, nil
}

// ═══════════════════════════════════════════════════════════════════
// FOLLOWS
// ═══════════════════════════════════════════════════════════════════

func Follow(db *sql.DB, followerID, followedID int) error {
	if followerID == followedID {
		return errors.New("impossible de se suivre soi-même")
	}

	var exists int
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE id = ?`, followedID).Scan(&exists)
	if exists == 0 {
		return errors.New("utilisateur introuvable")
	}

	_, err := db.Exec(
		`INSERT OR IGNORE INTO follows (follower_id, followed_id) VALUES (?, ?)`,
		followerID, followedID,
	)
	return err
}

func Unfollow(db *sql.DB, followerID, followedID int) error {
	_, err := db.Exec(
		`DELETE FROM follows WHERE follower_id = ? AND followed_id = ?`,
		followerID, followedID,
	)
	return err
}

func GetFollowing(db *sql.DB, userID int) ([]*models.User, error) {
	rows, err := db.Query(`
		SELECT u.id, u.nom, u.pseudo, u.email, u.photo_url, u.created_at
		FROM users u
		JOIN follows f ON f.followed_id = u.id
		WHERE f.follower_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUsers(rows)
}

func GetFollowers(db *sql.DB, userID int) ([]*models.User, error) {
	rows, err := db.Query(`
		SELECT u.id, u.nom, u.pseudo, u.email, u.photo_url, u.created_at
		FROM users u
		JOIN follows f ON f.follower_id = u.id
		WHERE f.followed_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUsers(rows)
}

func scanUsers(rows *sql.Rows) ([]*models.User, error) {
	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		rows.Scan(&u.ID, &u.Nom, &u.Pseudo, &u.Email, &u.PhotoURL, &u.CreatedAt)
		users = append(users, u)
	}
	return users, nil
}