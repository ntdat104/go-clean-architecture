package repo

import (
	"context"

	"github.com/ntdat104/go-clean-architecture/domain/model"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user *model.User) error
	UpdateWithTransaction(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
}

type UserCacheRepo interface {
	HealthCheck(ctx context.Context) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	Set(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error
	Invalidate(ctx context.Context) error
}
