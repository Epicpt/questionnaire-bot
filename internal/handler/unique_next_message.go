package handler

import (
	"fmt"
	"questionnaire-bot/internal/entity"
)

func (t *BotHandler) uniqueNextMessage(user *entity.User, message string) error {
	switch message {
	case combo:
		return t.comboMessage(user.TgID, user.ChatID)
	case attachments:
	//
	default:
		return fmt.Errorf("handler -> uniqueNextMessage -> wrong unique message: %s", message)
	}
	return nil
}

func (t *BotHandler) comboMessage(id int64, chatID int64) error {
	msg, advice, err := t.uc.GetComboMessage(id)
	if err != nil {
		return err
	}
	t.Send(chatID, msg.Text, nil)

	if advice != nil {
		t.Send(chatID, string(*advice), nil)
	}
	return nil
}
