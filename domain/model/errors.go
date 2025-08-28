package model

import (
	"github.com/ntdat104/go-clean-architecture/pkg/errors"
)

var (
	// user
	ErrInvalidUserID  = errors.New(errors.ErrorTypeValidation, "invalid user ID")
	ErrEmptyUserName  = errors.New(errors.ErrorTypeValidation, "user name cannot be empty")
	ErrEmptyUserEmail = errors.New(errors.ErrorTypeValidation, "user email cannot be empty")

	// example
	ErrExampleNotFound      = errors.New(errors.ErrorTypeNotFound, "example not found")
	ErrEmptyExampleName     = errors.New(errors.ErrorTypeValidation, "example name cannot be empty")
	ErrInvalidExampleID     = errors.New(errors.ErrorTypeValidation, "invalid example ID")
	ErrExampleNameTaken     = errors.New(errors.ErrorTypeConflict, "example name already taken")
	ErrExampleInvalidUpdate = errors.New(errors.ErrorTypeValidation, "invalid example update data")
	ErrExampleModified      = errors.New(errors.ErrorTypeConflict, "example modified by another process")
)
