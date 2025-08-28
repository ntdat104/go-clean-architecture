package model

import (
	"time"

	"github.com/ntdat104/go-clean-architecture/pkg/uuid"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Uuid      string    `json:"uuid" db:"uuid"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser(name, email string) (*User, error) {
	if name == "" {
		return nil, ErrEmptyUserName
	}
	if email == "" {
		return nil, ErrEmptyUserEmail
	}
	user := &User{
		Uuid:      uuid.NewShortUUID(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return user, nil
}

func (e *User) Validate() error {
	if e.Name == "" {
		return ErrEmptyUserName
	}
	if e.Email == "" {
		return ErrEmptyUserEmail
	}
	if e.ID < 0 {
		return ErrInvalidExampleID
	}
	return nil
}

func (e *User) Update(name string) error {
	if name == "" {
		return ErrEmptyUserName
	}
	e.Name = name
	e.UpdatedAt = time.Now()
	return nil
}
