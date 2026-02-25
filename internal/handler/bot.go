package handler

import (
	"fmt"
	"questionnaire-bot/internal/constantses"
	"questionnaire-bot/internal/entity"
	"time"
)

func remindUser(user *entity.User) {
	now := time.Now().UTC()
	user.RemindStage = constantses.Remind
	user.RemindAt = now.Add(constantses.RemindTime)
}

func (t *BotHandler) ProcessMessage(user *entity.User, text string) {

	if text == StartButton {
		user.CurrentStep = 0
		user.RemindStage = constantses.NotRemind
		user.IsCompleted = false // для повторного прохождения

		t.SendNextQuestion(user)
		return
	}

	//if user.IsCompleted {
	//	t.Send(user.ChatID, completedText, nil) // для повторного прохождения
	//	return
	//}

	if text == BackButton && user.CurrentStep >= 1 {
		user.CurrentStep--

		t.SendNextQuestion(user)
		return
	}

	var q *Question

	for _, answer := range NotifyAnswers {
		if answer.Text == text {
			q = &NotifyQuestion
		}
	}

	if q == nil {
		step := user.CurrentStep
		if step >= len(Questions) {
			t.FinishSurvey(user)
			return
		}

		q = &Questions[step]
	}

	if q.Key == "advice_2" {
		remindUser(user)
	}

	var err error
	var ans *Answer
	if q.Options != nil {
		ans, err = ValidateChoose(text, q.Options)
	} else if q.Validator != nil {
		err = q.Validator(text)
		if err == nil && q.InputAnswer != nil {
			ans = q.InputAnswer
			ans.Text = text
		}
	}
	if err != nil {
		t.Send(user.ChatID, fmt.Sprintf("%v", err), KeyboardFromOptions(*q, ShowBackButton(user.CurrentStep)))
		return
	}

	if ans != nil {
		answer := &entity.Answer{UserTgID: user.TgID, QuestionKey: q.Key, Step: user.CurrentStep, UserAnswer: ans.Text, Short: q.Short, TechName: ans.TechName}
		if err = t.uc.SaveAnswer(answer); err != nil {
			t.l.Err(err).Int64("userTgID", user.TgID).Str("text", text).Msg("failed to save answer")
			t.SendTo(t.ed.AdminID, adminPGErrorMessage())
			return
		}

		t.l.Info().Int64("tg id", answer.UserTgID).Str("username", user.Username).Str("question", answer.Short).Int("step", answer.Step).Str("answer", answer.UserAnswer).Msg("Answer success save")

		if ans.Actions != nil {
			if err = t.triggerMessage(user, ans); err != nil {
				t.SendTo(t.ed.AdminID, err.Error())
			}
		}
	}

	t.AdvanceStep(user)
}

func (t *BotHandler) ProcessContact(user *entity.User, phone string) {

	if user.IsCompleted {
		t.Send(user.ChatID, completedText, nil)
		return
	}
	q := Questions[user.CurrentStep]
	answer := &entity.Answer{UserTgID: user.TgID, QuestionKey: q.Key, Step: user.CurrentStep, UserAnswer: phone, Short: q.Short, TechName: "phone"}
	if err := t.uc.SaveAnswer(answer); err != nil {
		t.l.Err(err).Int64("userTgID", user.TgID).Str("phone", phone).Msg("failed to save answer")
		t.SendTo(t.ed.AdminID, adminPGErrorMessage())
		return
	}

	if err := t.alertManagers(user, constantses.ActionClientSentPhone); err != nil {
		t.SendTo(t.ed.AdminID, err.Error())
	}

	t.l.Info().Int64("tg id", answer.UserTgID).Str("username", user.Username).Str("question", answer.Short).Int("step", answer.Step).Str("answer", answer.UserAnswer).Msg("Answer success save")

	t.AdvanceStep(user)
}

func (t *BotHandler) AdvanceStep(user *entity.User) {
	user.CurrentStep++

	if user.CurrentStep >= len(Questions) {
		//t.FinishSurvey(user)
		return
	}

	t.SendNextQuestion(user)

}

func (t *BotHandler) SendNextQuestion(user *entity.User) {
	q := Questions[user.CurrentStep]
	t.Send(user.ChatID, q.Text, KeyboardFromOptions(q, ShowBackButton(user.CurrentStep)))

	if q.UniqueNextMessage != "" {
		if err := t.uniqueNextMessage(user, q.UniqueNextMessage); err != nil {
			t.SendTo(t.ed.AdminID, fmt.Sprintf("tgID: %d, name: %s, username: %s, error:%s", user.TgID, user.FirstName, user.Username, err.Error()))
		}
		t.AdvanceStep(user)
	}
}

func (t *BotHandler) FinishSurvey(user *entity.User) {
	user.IsCompleted = true

	t.Send(user.ChatID, finishText, nil)
}

func adminPGErrorMessage() string {
	return fmt.Sprintf("❗️Проблемы с БД❗️\nВозможно Postgres упал в %v.", time.Now().UTC().Format(time.RFC3339))
}

func ShowBackButton(step int) bool {
	if step < 6 && step > 0 {
		return true
	}
	return false
}
