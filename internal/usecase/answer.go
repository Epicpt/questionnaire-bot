package usecase

import "questionnaire-bot/internal/entity"

func (u *BotService) SaveAnswer(ans *entity.Answer) error {
	if err := u.repo.SaveAnswer(ans); err != nil {
		return err
	}
	return nil
}
