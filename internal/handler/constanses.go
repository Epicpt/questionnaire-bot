package handler

import "time"

const (
	remindTime  = time.Hour
	NotRemind   = 2 // dont send notification
	StartRemind = 0

	BackButton  = "⬅️ Назад"
	StartButton = "/start"

	requestPhone = "request_phone"

	// unique messages
	attachments = "attachments"
	combo       = "combo"

	// answer triggers
	alertManagers = "alertManagers"

	completedText = "Вы уже завершили анкету. Напишите /start, чтобы пройти заново."
	finishText    = "Спасибо! Вы прошли все вопросы. Менеджер свяжется с Вами в ближайшее время"
)
