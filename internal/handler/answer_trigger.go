package handler

import (
	"fmt"
	"questionnaire-bot/internal/constantses"
	"questionnaire-bot/internal/entity"
)

func (t *BotHandler) triggerMessage(user *entity.User, ans *Answer) error {
	switch ans.Trigger {
	case constantses.PhoneAlert:
		return t.alertManagers(user, ans.Trigger)
	case constantses.AppointmentAlert:
		if ans.TechName != "adv2a1" {
			return nil
		}
		return t.alertManagers(user, ans.Trigger)
	default:
		return fmt.Errorf("handler -> triggerMessage -> wrong trigger: %s", ans.Trigger)
	}
}

func (t *BotHandler) alertManagers(user *entity.User, trigger constantses.Trigger) error {
	msg, err := t.uc.GetManagerNotifyMessage(user, trigger)
	if err != nil {
		return err
	}
	for _, m := range t.ed.ManagerIDs {
		t.SendTo(m, *msg)
	}
	user.RemindStage = constantses.NotRemind
	return nil
}
