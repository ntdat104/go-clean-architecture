package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ntdat104/go-clean-architecture/internal/model"
	"github.com/ntdat104/go-clean-architecture/internal/repo"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) repo.UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `INSERT INTO users (name, email) VALUES (?, ?) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET name = ?, email = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = ?`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, err
	}
	return user, nil
}
