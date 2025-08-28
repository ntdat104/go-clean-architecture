package converter

import (
	"github.com/ntdat104/go-clean-architecture/application/dto"
	"github.com/ntdat104/go-clean-architecture/domain/model"
)

func ConvertUserToUserDto(input *model.User) *dto.User {
	return &dto.User{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	}
}
