package usecases

import (
	"errors"
	"regexp"
	"strings"

	"forum-diapason/internal/domain/entities"
	"forum-diapason/internal/domain/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Register(email, pseudo, password string) (*entities.User, error)
	Login(emailOrPseudo, password string) (*entities.User, error)
	GetProfile(userID int) (*entities.User, error)
	UpdateProfile(userID int, email, pseudo, photoURL string) error
	ChangePassword(userID int, oldPassword, newPassword string) error
}

type userUseCase struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (uc *userUseCase) Register(email, pseudo, password string) (*entities.User, error) {
	// Validate input
	if err := uc.validateEmail(email); err != nil {
		return nil, err
	}
	if err := uc.validatePseudo(pseudo); err != nil {
		return nil, err
	}
	if err := uc.validatePassword(password); err != nil {
		return nil, err
	}

	// Check if user already exists
	if _, err := uc.userRepo.GetByEmail(email); err == nil {
		return nil, errors.New("email already exists")
	}
	if _, err := uc.userRepo.GetByPseudo(pseudo); err == nil {
		return nil, errors.New("pseudo already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := entities.NewUser(email, pseudo, string(hashedPassword))
	err = uc.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) Login(emailOrPseudo, password string) (*entities.User, error) {
	var user *entities.User
	var err error

	// Try to find user by email or pseudo
	if strings.Contains(emailOrPseudo, "@") {
		user, err = uc.userRepo.GetByEmail(emailOrPseudo)
	} else {
		user, err = uc.userRepo.GetByPseudo(emailOrPseudo)
	}
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (uc *userUseCase) GetProfile(userID int) (*entities.User, error) {
	return uc.userRepo.GetByID(userID)
}

func (uc *userUseCase) UpdateProfile(userID int, email, pseudo, photoURL string) error {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if email != "" {
		if err := uc.validateEmail(email); err != nil {
			return err
		}
		user.Email = email
	}

	if pseudo != "" {
		if err := uc.validatePseudo(pseudo); err != nil {
			return err
		}
		user.Pseudo = pseudo
	}

	user.PhotoURL = photoURL
	return uc.userRepo.Update(user)
}

func (uc *userUseCase) ChangePassword(userID int, oldPassword, newPassword string) error {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}

	// Validate new password
	if err := uc.validatePassword(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return uc.userRepo.Update(user)
}

func (uc *userUseCase) validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func (uc *userUseCase) validatePseudo(pseudo string) error {
	if len(pseudo) < 3 || len(pseudo) > 20 {
		return errors.New("pseudo must be between 3 and 20 characters")
	}
	pseudoRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !pseudoRegex.MatchString(pseudo) {
		return errors.New("pseudo can only contain letters, numbers, underscores and hyphens")
	}
	return nil
}

func (uc *userUseCase) validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}
