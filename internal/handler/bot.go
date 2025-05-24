package handler

import (
	"fmt"
	"questionnaire-bot/internal/entity"
	"time"
)

type Question struct {
	Key             string
	Short           string
	Text            string
	Options         []string
	Validator       func(string) error
	SpecialKeyboard string
}

const (
	remindTime  = time.Hour
	NotRemind   = 2 // dont send notification
	StartRemind = 0

	BackButton  = "⬅️ Назад"
	StartButton = "/start"

	requestPhone = "request_phone"

	completedText = "Вы уже завершили анкету. Напишите /start, чтобы пройти заново."
	finishText    = "Спасибо! Вы прошли все вопросы. Менеджер свяжется с Вами в ближайшее время"
)

var (
	objectOpt     = []string{"Дом", "Квартира", "Коммерческая недвижимость"}
	renovationOpt = []string{"Капитальный (новостройка в бетоне)", "Капитальный (вторичка)", "Косметический (вторичка)"}
	equipmentOpt  = []string{"Да", "Нет"}
	connectionOpt = []string{"Телефонный звонок", "Написать в WhatsApp"}
)
var Questions = []Question{
	{
		Key: "name", Text: `Давайте познакомимся.
Уточните пожалуйста как Вас зовут?`,
		Validator: ValidateName,
		Short:     "Имя",
	},
	{
		Key: "object", Text: `Что у вас за объект?`,
		Options: objectOpt,
		Short:   "Вид объекта",
	},
	{
		Key: "area", Text: `Какая площадь объекта?`,
		Validator: ValidateArea,
		Short:     "Площадь",
	},
	{
		Key: "city", Text: `В каком городе находится объект?`,
		Validator: ValidateCity,
		Short:     "Город",
	},
	{
		Key: "renovation", Text: `Какой ремонт потребуется?`,
		Options: renovationOpt,
		Short:   "Вид ремонта",
	},
	{
		Key: "equipment", Text: `Нужна ли будет комплектация объекта во время ремонта?

Комплектация - это подбор мебели, светильников, декора, а также закупка всего необходимого для ремонта (сантехника, материалы, мебель, свтеильники и т.д.)`,
		Options: equipmentOpt,
		Short:   "Необходима комплектация",
	},
	{
		Key: "connection", Text: `Для уточнения дополнительных деталей проекта, рекомендуем связаться с нашим менеджером, чтобы подробно обсудить будущий дизайн-проект
	
	Какой способ связи для вас наиболее удобный?`,
		Options: connectionOpt,
		Short:   "Способ связи",
	},
	{
		Key: "phone", Text: `Оставьте свой номер телефона и наш менеджер свяжется с Вами в ближайшее время`,
		Validator:       ValidatePhone,
		SpecialKeyboard: requestPhone,
		Short:           "Номер телефона",
	},
}

func remindUser(user *entity.User) *entity.User {
	now := time.Now().UTC()
	user.RemindStage = StartRemind
	user.RemindAt = now.Add(remindTime)
	return user
}

func (t *BotHandler) ProcessMessage(user *entity.User, text string) {

	if text == StartButton {
		user.CurrentStep = 0
		user.IsCompleted = false // для повторного прохождения

		user = remindUser(user)

		t.SendNextQuestion(user)
		return
	}

	if user.IsCompleted {
		t.Send(user.ChatID, completedText, nil)
		return
	}

	if text == BackButton && user.CurrentStep >= 1 {
		user.CurrentStep--

		user = remindUser(user)

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
	if q.Options != nil {
		err = ValidateChoose(text, q.Options)
	} else {
		err = q.Validator(text)
	}
	if err != nil {
		t.Send(user.ChatID, fmt.Sprintf("%v", err), KeyboardFromOptions(q, user.CurrentStep > 0))
		return
	}

	answer := &entity.Answer{UserTgID: user.TgID, QuestionKey: q.Key, Step: step, UserAnswer: text, Short: q.Short}
	if err := t.u.SaveAnswer(answer); err != nil {
		t.l.Err(err).Msg("failed to save answer")
		t.SendToAdmin(adminMessage())
		return
	}

	t.l.Info().Int64("tg id", answer.UserTgID).Str("username", user.Username).Str("question", answer.Short).Int("step", answer.Step).Str("answer", answer.UserAnswer).Msg("Answer success save")

	t.AdvanceStep(user)
}

func (t *BotHandler) ProcessContact(user *entity.User, phone string) {

	if user.IsCompleted {
		t.Send(user.ChatID, completedText, nil)
		return
	}
	q := Questions[user.CurrentStep]
	answer := &entity.Answer{UserTgID: user.TgID, QuestionKey: q.Key, Step: user.CurrentStep, UserAnswer: phone, Short: q.Short}
	if err := t.u.SaveAnswer(answer); err != nil {
		t.l.Err(err).Msg("failed to save answer")
		t.SendToAdmin(adminMessage())
		return
	}

	t.l.Info().Int64("tg id", answer.UserTgID).Str("username", user.Username).Str("question", answer.Short).Int("step", answer.Step).Str("answer", answer.UserAnswer).Msg("Answer success save")

	t.AdvanceStep(user)
}

func (t *BotHandler) AdvanceStep(user *entity.User) {
	user.CurrentStep++

	user = remindUser(user)

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

	if err := t.u.CreateEmail(user); err != nil {
		t.l.Err(err).Msg("failed to create email")
		t.SendToAdmin(adminMessage())
		return
	}
	user.EmailSentCnt++
	t.Send(user.ChatID, finishText, nil)
}

func adminMessage() string {
	return fmt.Sprintf("❗️Проблемы с БД❗️\nВозможно Postgres упал в %v.", time.Now().UTC().Format(time.RFC3339))
}
