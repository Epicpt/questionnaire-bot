package repository

import (
	"questionnaire-bot/internal/entity"
	"questionnaire-bot/pkg/postgres"
)

type Repository interface {
	GetUser(int64) (*entity.User, error)
	SaveUser(*entity.User) error
	SaveAnswer(*entity.Answer) error
	GetAnswersUser(int64) ([]entity.Answer, error)
	SaveEmail(*entity.Email) error
	GetEmailsByStatus(string) ([]entity.Email, error)
	UpdateEmailStatus(*entity.Email, string) error
	GetUsersForNotify() ([]entity.User, error)
}

type BotRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *BotRepo {
	return &BotRepo{pg}
}
