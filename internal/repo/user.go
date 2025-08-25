package repo

import (
	"context"

	"github.com/ntdat104/go-clean-architecture/internal/model"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
}
