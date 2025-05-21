package usecase

import (
	"questionnaire-bot/internal/entity"
)

func (u *BotService) GetUsersForNotify() ([]entity.User, error) {
	return u.repo.GetUsersForNotify()
}
