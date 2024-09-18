package repo

import (
	"errors"
	"github.com/Frozelo/startupFeed/internal/models"
)

type InMemoryRepository struct {
	users map[string]*models.User
}

func NewInMemoryRepository(users map[string]*models.User) *InMemoryRepository {
	return &InMemoryRepository{users: users}
}

func (r *InMemoryRepository) GetUser(username string) (*models.User, error) {
	user, ok := r.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryRepository) CreateUser(user *models.User) error {
	_, ok := r.users[user.Username]
	if ok {
		return errors.New("user already exists")
	}

	r.users[user.Username] = user

	return nil
}

func (r *InMemoryRepository) AddComment(userId int, comment *models.Comment) error {
	comment.UserId = userId

	return nil
}
