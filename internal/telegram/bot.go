package telegram

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
	notRemind   = 2 // dont send notification
	startRemind = 0

	backButton  = "⬅️ Назад"
	startButton = "/start"

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
		Validator: validateName,
		Short:     "Имя",
	},
	{
		Key: "object", Text: `Что у вас за объект?`,
		Options:   objectOpt,
		Validator: validateObject,
		Short:     "Вид объекта",
	},
	{
		Key: "area", Text: `Какая площадь объекта?`,
		Validator: validateArea,
		Short:     "Площадь",
	},
	{
		Key: "city", Text: `В каком городе находится объект?`,
		Validator: validateCity,
		Short:     "Город",
	},
	{
		Key: "renovation", Text: `Какой ремонт потребуется?`,
		Options:   renovationOpt,
		Validator: validateRenovation,
		Short:     "Вид ремонта",
	},
	{
		Key: "equipment", Text: `Нужна ли будет комплектация объекта во время ремонта?

Комплектация - это подбор мебели, светильников, декора, а также закупка всего необходимого для ремонта (сантехника, материалы, мебель, свтеильники и т.д.)`,
		Options:   equipmentOpt,
		Validator: validateEquipment,
		Short:     "Необходима комплектация",
	},
	{
		Key: "connection", Text: `Для уточнения дополнительных деталей проекта, рекомендуем связаться с нашим менеджером, чтобы подробно обсудить будущий дизайн-проект
	
	Какой способ связи для вас наиболее удобный?`,
		Options:   connectionOpt,
		Validator: validateConnection,
		Short:     "Способ связи",
	},
	{
		Key: "phone", Text: `Оставьте свой номер телефона и наш менеджер свяжется с Вами в ближайшее время`,
		Validator:       validatePhone,
		SpecialKeyboard: requestPhone,
		Short:           "Номер телефона",
	},
}

func (t *Telegram) processMessage(user *entity.User, text string) {
	now := time.Now().UTC()

	if text == startButton {
		user.CurrentStep = 0
		user.IsCompleted = false // для повторного прохождения
		user.RemindStage = startRemind
		user.RemindAt = now.Add(remindTime)

		t.sendNextQuestion(user)
		return
	}

	if user.IsCompleted {
		t.Send(user.ChatID, completedText, nil)
		return
	}

	if text == backButton && user.CurrentStep >= 1 {
		user.CurrentStep--
		user.RemindAt = now.Add(remindTime)

		t.sendNextQuestion(user)
		return
	}

	step := user.CurrentStep
	if step >= len(Questions) {
		t.finishSurvey(user)
		return
	}

	q := Questions[step]
	if q.Validator != nil {
		if err := q.Validator(text); err != nil {
			t.Send(user.ChatID, fmt.Sprintf("%v", err), keyboardFromOptions(q, user.CurrentStep > 0))
			return
		}
	}
	answer := &entity.Answer{UserTgID: user.TgID, QuestionKey: q.Key, Step: step, UserAnswer: text, Short: q.Short}
	t.u.SaveAnswer(answer)

	t.l.Info().Int64("tg id", answer.UserTgID).Str("username", user.Username).Str("question", answer.Short).Int("step", answer.Step).Str("answer", answer.UserAnswer).Msg("Answer success save")

	t.advanceStep(user)
}

func (t *Telegram) processContact(user *entity.User, phone string) {

	if user.IsCompleted {
		t.Send(user.ChatID, completedText, nil)
		return
	}
	q := Questions[user.CurrentStep]
	answer := &entity.Answer{UserTgID: user.TgID, QuestionKey: q.Key, Step: user.CurrentStep, UserAnswer: phone, Short: q.Short}
	t.u.SaveAnswer(answer)

	t.l.Info().Int64("tg id", answer.UserTgID).Str("username", user.Username).Str("question", answer.Short).Int("step", answer.Step).Str("answer", answer.UserAnswer).Msg("Answer success save")

	t.advanceStep(user)
}

func (t *Telegram) advanceStep(user *entity.User) {
	now := time.Now().UTC()
	user.CurrentStep++
	user.RemindAt = now.Add(remindTime)

	if user.CurrentStep >= len(Questions) {
		t.finishSurvey(user)
		return
	}

	t.sendNextQuestion(user)
}

func (t *Telegram) sendNextQuestion(user *entity.User) {
	q := Questions[user.CurrentStep]
	t.Send(user.ChatID, q.Text, keyboardFromOptions(q, user.CurrentStep > 0))
}

func (t *Telegram) finishSurvey(user *entity.User) {
	user.IsCompleted = true
	user.RemindStage = notRemind

	if err := t.u.CreateEmail(user); err != nil {
		t.l.Err(err).Msg("failed to create email")
	}
	t.Send(user.ChatID, finishText, nil)
}
