package entity

import "time"

type (
	User struct {
		TgID           int64
		ChatID         int64
		FirstName      string
		LastName       string
		Username       string
		CreatedAt      time.Time
		UpdatedAt      time.Time
		RemindStage    int
		RemindAt       time.Time
		IsCompleted    bool
		CurrentStep    int
		MaxStepReached int
	}

	Answer struct {
		ID          int
		UserTgID    int64
		QuestionKey string
		Short       string
		Step        int
		UserAnswer  string
		CreatedAt   time.Time
	}

	Email struct {
		ID        int
		UserTgID  int64
		Body      string
		Status    string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func New(tgID, chatID int64, firstName, lastName, userName string) *User {
	return &User{TgID: tgID, ChatID: chatID, FirstName: firstName, LastName: lastName, Username: userName}
}
