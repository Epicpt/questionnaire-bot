package usecase

import (
	"fmt"
	"questionnaire-bot/internal/constantses"
	"questionnaire-bot/internal/entity"
	"strings"
)

func (u *BotService) GetUsersForNotify() ([]entity.User, error) {
	return u.repo.GetUsersForNotify()
}

func (u *BotService) GetManagerNotifyMessage(user *entity.User, trigger constantses.Trigger) (*string, error) {
	answers, err := u.repo.GetAnswersUser(user.TgID)
	if err != nil {
		return nil, fmt.Errorf("BotService -> GetManagerNotifyMessage -> repo.GetAnswersUser: %w", err)
	}

	var sb strings.Builder
	switch trigger {
	case constantses.PhoneAlert:
		sb.WriteString(fmt.Sprintf("Клиент отправил телефон для получения материалов. Клиент пока не дошел до записи на консультацию: \n\n---Внутренние данные телеграмм---\n"))
	case constantses.AppointmentAlert:
		sb.WriteString(fmt.Sprintf("Клиент отправил запрос на запись на консультацию: \n\n---Внутренние данные телеграмм---\n"))
	}

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

	res := sb.String()

	return &res, nil
}
