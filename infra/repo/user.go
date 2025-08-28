package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ntdat104/go-clean-architecture/domain/model"
	"github.com/ntdat104/go-clean-architecture/domain/repo"
	"github.com/ntdat104/go-clean-architecture/infra/repository"
)

// userRepo implements the repo.UserRepo interface using sqlx.
type userRepo struct {
	db *sqlx.DB
}

// NewUserRepo creates a new UserRepo with the provided sqlx.DB.
func NewUserRepo(client *repository.Client) repo.UserRepo {
	return &userRepo{db: client.MySQL}
}

// Create inserts a new user into the database and returns the created user with its new ID.
func (r *userRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (uuid, name, email, created_at, updated_at) 
		VALUES (:uuid, :name, :email, :created_at, :updated_at) 
		RETURNING id
	`
	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user.ID); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("failed to retrieve ID after insert")
	}
	return user, nil
}

// Delete removes a user from the database by ID.
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

// Update modifies an existing user in the database.
func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users 
		SET uuid = :uuid, name = :name, email = :email, updated_at = :updated_at 
		WHERE id = :id
	`
	result, err := r.db.NamedExecContext(ctx, query, user)
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

// Update modifies a user in a transaction.
func (r *userRepo) UpdateWithTransaction(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
		UPDATE users 
		SET uuid = :uuid, name = :name, email = :email, updated_at = :updated_at 
		WHERE id = :id
	`
	result, err := tx.NamedExecContext(ctx, query, user)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows == 0 {
		tx.Rollback()
		return errors.New("user not found")
	}

	return tx.Commit()
}

// GetByID retrieves a user from the database by ID.
func (r *userRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `SELECT id, uuid, name, email, created_at, updated_at FROM users WHERE id = ?`
	user := &model.User{}
	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
