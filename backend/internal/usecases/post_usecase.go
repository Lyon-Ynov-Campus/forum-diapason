package usecases

import (
	"errors"
	"strings"

	"forum-diapason/internal/entities"
	"forum-diapason/internal/repositories"
)

type PostUseCase interface {
	CreatePost(authorID int, title, content string) (*entities.Post, error)
	GetPost(id int) (*entities.Post, error)
	GetPosts(limit, offset int) ([]*entities.Post, error)
	GetUserPosts(userID, limit, offset int) ([]*entities.Post, error)
	UpdatePost(postID, authorID int, title, content string) error
	DeletePost(postID, authorID int) error
	LikePost(postID, userID int) error
	UnlikePost(postID, userID int) error
}

type postUseCase struct {
	postRepo repositories.PostRepository
	userRepo repositories.UserRepository
}

func NewPostUseCase(postRepo repositories.PostRepository, userRepo repositories.UserRepository) PostUseCase {
	return &postUseCase{postRepo: postRepo, userRepo: userRepo}
}

func (uc *postUseCase) CreatePost(authorID int, title, content string) (*entities.Post, error) {
	// Validate input
	if err := uc.validateTitle(title); err != nil {
		return nil, err
	}
	if err := uc.validateContent(content); err != nil {
		return nil, err
	}

	// Get author pseudo
	author, err := uc.userRepo.GetByID(authorID)
	if err != nil {
		return nil, err
	}

	// Create post
	post := &entities.Post{
		AuthorID:     authorID,
		AuthorPseudo: author.Pseudo,
		Title:        strings.TrimSpace(title),
		Content:      strings.TrimSpace(content),
		LikeCount:    0,
		CommentCount: 0,
	}
	err = uc.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (uc *postUseCase) GetPost(id int) (*entities.Post, error) {
	return uc.postRepo.GetByID(id)
}

func (uc *postUseCase) GetPosts(limit, offset int) ([]*entities.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.postRepo.GetAll(limit, offset)
}

func (uc *postUseCase) GetUserPosts(userID, limit, offset int) ([]*entities.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.postRepo.GetByAuthorID(userID, limit, offset)
}

func (uc *postUseCase) UpdatePost(postID, authorID int, title, content string) error {
	// Get existing post
	post, err := uc.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	// Check ownership
	if post.AuthorID != authorID {
		return errors.New("unauthorized")
	}

	// Validate input
	if title != "" {
		if err := uc.validateTitle(title); err != nil {
			return err
		}
		post.Title = strings.TrimSpace(title)
	}
	if content != "" {
		if err := uc.validateContent(content); err != nil {
			return err
		}
		post.Content = strings.TrimSpace(content)
	}

	return uc.postRepo.Update(post)
}

func (uc *postUseCase) DeletePost(postID, authorID int) error {
	// Get existing post
	post, err := uc.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	// Check ownership
	if post.AuthorID != authorID {
		return errors.New("unauthorized")
	}

	return uc.postRepo.Delete(postID)
}

func (uc *postUseCase) LikePost(postID, userID int) error {
	// Check if post exists
	_, err := uc.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	return uc.postRepo.Like(postID, userID)
}

func (uc *postUseCase) UnlikePost(postID, userID int) error {
	// Check if post exists
	_, err := uc.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	return uc.postRepo.Unlike(postID, userID)
}

func (uc *postUseCase) validateTitle(title string) error {
	title = strings.TrimSpace(title)
	if len(title) < 5 || len(title) > 200 {
		return errors.New("title must be between 5 and 200 characters")
	}
	return nil
}

func (uc *postUseCase) validateContent(content string) error {
	content = strings.TrimSpace(content)
	if len(content) < 10 || len(content) > 10000 {
		return errors.New("content must be between 10 and 10000 characters")
	}
	return nil
}
