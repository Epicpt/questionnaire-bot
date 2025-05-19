package usecase

import (
	"errors"
	"questionnaire-bot/internal/entity"
)

var ErrUserNotFound = errors.New("user not found")

func (u *BotService) SaveUser(user *entity.User) error {
	if err := u.repo.SaveUser(user); err != nil {
		return err
	}
	return nil
}

func (u *BotService) GetUser(userID int64) (*entity.User, error) {
	user, err := u.repo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
