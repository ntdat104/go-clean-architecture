package service

import (
	"time"

	"github.com/ntdat104/go-clean-architecture/application/dto"
)

type SystemService interface {
	GetTime() *dto.SystemTime
}

type systemService struct{}

func NewSystemService() SystemService {
	return &systemService{}
}

func (s *systemService) GetTime() *dto.SystemTime {
	return &dto.SystemTime{
		ServerTime: time.Now().UnixMilli(),
	}
}
