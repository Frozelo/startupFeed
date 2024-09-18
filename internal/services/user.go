package services

import (
	"errors"
	"github.com/Frozelo/startupFeed/internal/models"
	"github.com/Frozelo/startupFeed/internal/repo"
)

type MemoryRepo interface {
	GetUser(username string) (*models.User, error)
	CreateUser(user *models.User) error
	AddComment(userId int, comment *models.Comment) error
}

type UserService struct {
	repo *repo.InMemoryRepository
}

func NewUserService(repo *repo.InMemoryRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(username string) (*models.User, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	if user.Username == "" {
		return errors.New("invalid user data")
	}
	return s.repo.CreateUser(user)
}

func (s *UserService) AddComment(userId int, comment *models.Comment) error {
	return s.repo.AddComment(userId, comment)
}
