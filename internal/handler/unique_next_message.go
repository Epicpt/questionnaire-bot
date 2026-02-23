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
		return t.attachmentsMessage(user.TgID, user.ChatID)
	case skip:
		return nil
	default:
		return fmt.Errorf("handler -> uniqueNextMessage -> wrong unique message: %s", message)
	}
}

func (t *BotHandler) comboMessage(id int64, chatID int64) error {
	data, advice, err := t.uc.GetComboMessage(id)
	if err != nil {
		return err
	}
	t.Send(chatID, data.Text, nil)

	if advice != nil {
		t.Send(chatID, string(*advice), nil)
	}
	return nil
}

func (t *BotHandler) attachmentsMessage(id int64, chatID int64) error {
	data, _, err := t.uc.GetComboMessage(id)
	if err != nil {
		return err
	}

	t.Send(chatID, data.AttachmentText, nil)
	t.SendDocs(chatID, data.FilePath)
	return nil
}
