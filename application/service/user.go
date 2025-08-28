package service

import (
	"context"
	"errors"

	"github.com/ntdat104/go-clean-architecture/domain/model"
	"github.com/ntdat104/go-clean-architecture/domain/repo"
)

var (
	ErrInvalidUser = errors.New("invalid user data")
	ErrNotFound    = errors.New("user not found")
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	repo repo.UserRepo
}

func NewUserService(r repo.UserRepo) UserService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// basic validation
	if user.Name == "" || user.Email == "" {
		return nil, ErrInvalidUser
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	if user.ID == 0 || user.Name == "" || user.Email == "" {
		return ErrInvalidUser
	}
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	if id == 0 {
		return ErrInvalidUser
	}
	return s.repo.Delete(ctx, id)
}
