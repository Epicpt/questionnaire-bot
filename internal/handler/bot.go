package handler

import (
	"fmt"
	"questionnaire-bot/internal/entity"
	"time"
)

func remindUser(user *entity.User) {
	now := time.Now().UTC()
	user.RemindStage = StartRemind
	user.RemindAt = now.Add(remindTime)
}

func (t *BotHandler) ProcessMessage(user *entity.User, text string) {

	if text == StartButton {
		user.CurrentStep = 0
		user.IsCompleted = false // для повторного прохождения

		remindUser(user)

		t.SendNextQuestion(user)
		return
	}

	if user.IsCompleted {
		t.Send(user.ChatID, completedText, nil)
		return
	}

	if text == BackButton && user.CurrentStep >= 1 {
		user.CurrentStep--

		remindUser(user)

		t.SendNextQuestion(user)
		return
	}

	step := user.CurrentStep
	if step >= len(Questions) {
		t.FinishSurvey(user)
		return
	}

	q := Questions[step]

	var err error
	ans := new(Answer)
	if q.Options != nil {
		ans, err = ValidateChoose(text, q.Options)
	} else {
		err = q.Validator(text)
	}
	if err != nil {
		t.Send(user.ChatID, fmt.Sprintf("%v", err), KeyboardFromOptions(q, user.CurrentStep > 0))
		return
	}

	if ans != nil {
		answer := &entity.Answer{UserTgID: user.TgID, QuestionKey: q.Key, Step: step, UserAnswer: ans.Text, Short: q.Short, TechName: ans.TechName}
		if err := t.uc.SaveAnswer(answer); err != nil {
			t.l.Err(err).Msg("failed to save answer")
			t.SendTo(t.ed.AdminID, adminMessage())
			return
		}

		t.l.Info().Int64("tg id", answer.UserTgID).Str("username", user.Username).Str("question", answer.Short).Int("step", answer.Step).Str("answer", answer.UserAnswer).Msg("Answer success save")

		if ans.Trigger != "" {
			// todo: логика отправки уведомления менеджерам
		}
	}

	if q.UniqueNextMessage != "" {
		// todo: функция func(string) error с проверками на какое сообщение со своей логикой
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
		t.l.Err(err).Msg("failed to save answer")
		t.SendTo(t.ed.AdminID, adminMessage())
		return
	}

	t.l.Info().Int64("tg id", answer.UserTgID).Str("username", user.Username).Str("question", answer.Short).Int("step", answer.Step).Str("answer", answer.UserAnswer).Msg("Answer success save")

	if q.UniqueNextMessage != "" {
		// todo: функция func(string) error с проверками на какое сообщение со своей логикой
	}

	t.AdvanceStep(user)
}

func (t *BotHandler) AdvanceStep(user *entity.User) {
	user.CurrentStep++

	remindUser(user)

	if user.CurrentStep >= len(Questions) {
		t.FinishSurvey(user)
		return
	}

	t.SendNextQuestion(user)
}

func (t *BotHandler) SendNextQuestion(user *entity.User) {
	q := Questions[user.CurrentStep]
	t.Send(user.ChatID, q.Text, KeyboardFromOptions(q, user.CurrentStep > 0))
}

func (t *BotHandler) FinishSurvey(user *entity.User) {
	user.IsCompleted = true
	user.RemindStage = NotRemind

	if err := t.uc.CreateEmail(user); err != nil {
		t.l.Err(err).Msg("failed to create email")
		t.SendTo(t.ed.AdminID, adminMessage())
		return
	}
	user.EmailSentCnt++
	t.Send(user.ChatID, finishText, nil)
}

func adminMessage() string {
	return fmt.Sprintf("❗️Проблемы с БД❗️\nВозможно Postgres упал в %v.", time.Now().UTC().Format(time.RFC3339))
}
