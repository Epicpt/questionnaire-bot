package usecase

import (
	"questionnaire-bot/internal/constantses"
	"questionnaire-bot/internal/entity"
	"questionnaire-bot/internal/messages"
	"questionnaire-bot/internal/repository"
	"questionnaire-bot/pkg/smtp"
)

type Usecase interface {
	GetUser(int64) (*entity.User, error)
	SaveUser(*entity.User) error
	SaveAnswer(*entity.Answer) error
	CreateEmail(*entity.User) error
	GetEmailsByStatus(string) ([]entity.Email, error)
	SendEmail(*entity.Email) error
	UpdateEmailStatus(*entity.Email, string) error
	GetUsersForNotify() ([]entity.User, error)
	GetComboMessage(id int64) (*messages.Combo, *messages.PersonalAdvice, error)
	GetManagerNotifyMessage(user *entity.User, trigger constantses.Trigger) (*string, error)
}

type BotService struct {
	repo repository.Repository
	smtp smtp.EmailClient
}

func New(repo repository.Repository, smtp smtp.EmailClient) *BotService {
	return &BotService{repo: repo, smtp: smtp}
}
