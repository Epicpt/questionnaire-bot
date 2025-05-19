package usecase

import (
	"questionnaire-bot/internal/entity"
	"questionnaire-bot/internal/repository"
)

type Usecase interface {
	GetUser(int64) (*entity.User, error)
	SaveUser(*entity.User) error
	SaveAnswer(*entity.Answer) error
	CreateEmail(*entity.User) error
}

type BotService struct {
	repo repository.Repository
}

func New(repo repository.Repository) *BotService {
	return &BotService{repo: repo}
}
