package repository

import (
	"questionnaire-bot/internal/entity"
	"questionnaire-bot/pkg/postgres"
)

type Repository interface {
	GetUser(int64) (*entity.User, error)
	SaveUser(*entity.User) error
}

type BotRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *BotRepo {
	return &BotRepo{pg}
}
