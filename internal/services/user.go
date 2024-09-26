package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/Frozelo/startupFeed/internal/dto"
	"github.com/Frozelo/startupFeed/internal/models"
)

const (
	DefaultRole = "default"
	AdminRole   = "admin"
)

type UserRepo interface {
	Create(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetAuthors(ctx context.Context, projectId int64) ([]int64, error)
}

type UserService struct {
	userRepo UserRepo
}

func NewUserSerice(userRepo UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (us *UserService) Register(
	ctx context.Context,
	userDTO *dto.CreateUserDTO,
) error {
	var newUser *models.User
	passwordHash, err := passwordHash(userDTO.Password)
	if err != nil {
		return err
	}

	newUser = &models.User{
		Username:     userDTO.Username,
		Email:        userDTO.Email,
		PasswordHash: passwordHash,
		Role:         DefaultRole,
	}

	if err := us.userRepo.Create(ctx, newUser); err != nil {
		return err
	}
	return nil
}

func (us *UserService) Login(
	ctx context.Context,
	loginUserDTO *dto.LoginUserDTO,
) (*models.User, error) {
	user, err := us.userRepo.FindUserByEmail(ctx, loginUserDTO.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginUserDTO.Password)); err != nil {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	stringHash := string(hash)

	return stringHash, nil
}
