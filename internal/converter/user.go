package converter

import (
	"github.com/ntdat104/go-clean-architecture/internal/dto"
	"github.com/ntdat104/go-clean-architecture/internal/model"
)

func ConvertUserToUserDto(input *model.User) *dto.User {
	return &dto.User{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	}
}
