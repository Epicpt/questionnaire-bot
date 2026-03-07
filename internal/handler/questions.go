package handler

import (
	"questionnaire-bot/internal/constantses"
	"questionnaire-bot/internal/entity"
)

type Question struct {
	Key               string
	Short             string
	Text              string
	Options           []Answer
	Validator         func(string) error
	SpecialKeyboard   string
	UniqueNextMessage string
	InputAnswer       *Answer
	Photo             *entity.Photo
}

var Questions = []Question{
	{
		Key: "begin", Text: `🏠 ✨<b>Приветствую вас новосел и поздравляю с покупкой новой прекрасной квартиры!</b> 

Я бот <b>Solo Design Studio</b>, помогу провести быструю диагностику вашей <b>готовности к ремонту</b>

<i>Ответьте на 5 коротких вопросов — в конце вы получите:</i>

  ✅ <b>Отчет</b> и персональную <b>карту рисков</b> — что может пойти не так именно в вашем ремонте

  ✅ Полезные <b>чек-листы</b> для старта в ремонте - чек-лист расходов на ремонт, чек-лист выбора строительной бригады и другие 

<i>Это бесплатно и займет не более 3 минут</i>
 
<b>Начнем?</b>`,
		Options: beginOpt,
		Short:   "Приветственное сообщение",
		Photo: &entity.Photo{
			Paths: []string{"materials/photo_start.jpg"},
			Type:  entity.Jpeg,
		},
	},
	{
		Key: "q_1", Text: `❓<i>Вопрос 1/5</i> 

Ваш <b>идеальный ремонт</b> — это...`,
		Options: q1Opt,
		Short:   "Идеальный ремонт",
	},
	{
		Key: "q_2", Text: `❓<i>Вопрос 2/5</i>

Что <b>пугает</b> вас в предстоящем ремонте больше всего?`,
		Options: q2Opt,
		Short:   "Пугает в ремонте",
	},
	{
		Key: "q_3", Text: `❓<i>Вопрос 3/5</i>

Когда ремонт будет закончен, что для вас станет <b>главным показателем успеха?</b>`,
		Options: q3Opt,
		Short:   "Показатель успеха",
	},
	{
		Key: "q_4", Text: `❓<i>Вопрос 4/5</i>

Как вы планируете принимать <b>ключевые решения</b> по ремонту?`,
		Options: q4Opt,
		Short:   "Ключевые решения",
	},
	{
		Key: "q_5", Text: `❓<i>Вопрос 5/5</i>

На какой вы сейчас <b>стадии?</b>`,
		Options: q5Opt,
		Short:   "Стадия ремонта",
	},
	{
		Key:               "download_sobaka_ulibaka_gif",
		Short:             "Гифка загрузки",
		UniqueNextMessage: skip,
		Photo: &entity.Photo{
			Paths: []string{"materials/download_gif.mp4"},
			Type:  entity.Animation,
		},
	},
	{
		Key: "ready_combo", Text: `🎉 <i>Отлично! Ваша диагностика завершена!</i>

<b>Отчет</b> о вашей готовности к ремонту и <b>персональная карта рисков — ГОТОВЫ!</b>`,
		Short:             "Опрос завершен",
		UniqueNextMessage: combo,
	},

	{
		Key: "phone", Text: `🎁 Не уходите сразу, у нас еще остались <b>подарки</b> для вас!

Чтобы <u>получить чек-листы и гайды</u>, оставьте свой контакт — номер телефона или Telegram-никнейм`,
		Validator:       ValidatePhone,
		SpecialKeyboard: requestPhone,
		Short:           "Номер телефона",
		InputAnswer: &Answer{
			TechName: "phone",
			Actions:  []constantses.Action{constantses.ActionClientSentPhone},
		},
		Photo: &entity.Photo{
			Paths: []string{"materials/photo_gift_phone.jpg"},
			Type:  entity.Jpeg,
		},
	},
	{
		Key: "materials", Text: `Ваши материалы уже здесь 👇🏼`,
		Short:             "Отправили материалы",
		UniqueNextMessage: attachments,
	},
	{
		Key: "advice_1",
		Text: `💡 <b>СОВЕТ</b>: перед тем, как вы пойдете изучать материалы и готовиться к ремонту - обдумайте план действий и как вы видите будущий ремонт

Узнайте, как прошел ремонт у друзей и знакомых. Это даст больше информации для сравнения

📢 <i>Больше практических советов — у нас в Telegram-канале <a href="https://t.me/solo_design_studio">Solo Design Studio</a>
Подписывайтесь!</i>`,
		UniqueNextMessage: skip,
	},
	{
		Key:   "advice_2",
		Short: "Запись на консультацию?",
		Text: `❓Остались вопросы, приглашаем вас на <b>бесплатную консультацию!</b>

<b>Мы поможем разобраться:</b>
✅ как лучше спланировать пространство,
✅ как определить бюджет на ремонт,
✅ какие цвета подойдут под стиль, который вам нравится

📝 Записаться на бесплатную консультацию можно, написав 👉 @solo_ds`,
		Options: advice2Opt,
		Photo: &entity.Photo{
			Paths: []string{"materials/photo_advice_2.jpg"},
			Type:  entity.Jpeg,
		},
	},
}

var NotifyQuestion = Question{
	Key: "n_q_1",
	Text: `Хотите понять с чего начать ремонт в квартире? 
Бесплатная консультация от Solo Design Studio - отличная возможность чтобы:
- Разобраться как спланировать пространство, чтобы все было удобно функционально
- Уточнить все вопросы по цветам и материалам для ремонта - что выбрать в многообразии ремонтного мира
- Спланировать бюджет на ремонт и как учесть все риски

Записать вас на консультацию?`,
	Options: NotifyAnswers,
	Short:   "Напоминание о записи на консультацию",
}
