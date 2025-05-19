package usecase

import (
	"fmt"
	"questionnaire-bot/internal/entity"
	"strings"
)

const pendingStatus = "pending"

func (u *BotService) CreateEmail(user *entity.User) error {
	answers, err := u.repo.GetAnswersUser(user.TgID)
	if err != nil {
		return fmt.Errorf("fetching answers: %w", err)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Новая заявка: \n\n---Внутренние данные телеграмм---"))

	if user.FirstName != "" {
		sb.WriteString(fmt.Sprintf("* Имя: %s\n", user.FirstName))
	}
	if user.LastName != "" {
		sb.WriteString(fmt.Sprintf("* Фамилия: %s\n", user.LastName))
	}
	if user.Username != "" {
		sb.WriteString(fmt.Sprintf("* Username: %s\n", user.Username))
	}

	sb.WriteString("------------------------------------------\nОтветы по анкете:\n")
	for _, a := range answers {
		sb.WriteString(fmt.Sprintf("* %s: %s\n", a.Short, a.UserAnswer))
	}

	email := &entity.Email{
		UserTgID: user.TgID,
		Body:     sb.String(),
		Status:   pendingStatus,
	}

	if err := u.repo.SaveEmail(email); err != nil {
		return fmt.Errorf("failed save email: %w", err)
	}

	return nil
}
