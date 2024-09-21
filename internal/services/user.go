package services

import (
	"context"

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

	// TODO userRepo logic: exsiting check
	//===================================

	// TODO hashing password logic
	//===================================

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

	// TODO implemet userRepo logic: save user data
	// =================================
	if err := us.userRepo.Create(ctx, newUser); err != nil {
		return err
	}
	return nil
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
