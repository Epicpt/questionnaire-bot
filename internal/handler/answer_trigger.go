package handler

import (
	"fmt"
	"questionnaire-bot/internal/constantses"
	"questionnaire-bot/internal/entity"
	"questionnaire-bot/internal/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *BotHandler) triggerMessage(user *entity.User, ans *Answer) error {
	for _, action := range ans.Actions {
		switch action {
		case constantses.ActionClientSentPhone:
			if err := t.alertManagers(user, action); err != nil {
				return err
			}
		case constantses.ActionClientSentAppointment:
			if err := t.alertManagers(user, action); err != nil {
				return err
			}
		case constantses.ActionSendBookingMessage:
			t.Send(user.ChatID, messages.Booking, tgbotapi.NewRemoveKeyboard(true))
		case constantses.ActionSendDeclineMessage:
			t.Send(user.ChatID, messages.Decline, tgbotapi.NewRemoveKeyboard(true))
		default:
			return fmt.Errorf("handler -> triggerMessage -> wrong action: %s", ans.Actions)
		}
	}
	return nil
}

func (t *BotHandler) alertManagers(user *entity.User, action constantses.Action) error {
	msg, err := t.uc.GetManagerNotifyMessage(user, action)
	if err != nil {
		return err
	}
	for _, m := range t.ed.ManagerIDs {
		t.SendTo(m, *msg)
	}
	user.RemindStage = constantses.NotRemind
	return nil
}
