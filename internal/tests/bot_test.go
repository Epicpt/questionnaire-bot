package tests

import (
	"questionnaire-bot/internal/entity"
	"questionnaire-bot/internal/handler"
	"questionnaire-bot/internal/mocks"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

func TestBotProcessMessage(t *testing.T) {
	testCases := []struct {
		desc       string
		user       *entity.User
		text       string
		setupMocks func(telegram *mocks.MockTelegram, usecase *mocks.MockUsecase)
		exeptUser  *entity.User
	}{
		{
			desc: "Start command resets progress and sends next question",
			user: &entity.User{
				CurrentStep: 10,
				IsCompleted: true,
				RemindStage: 2,
			},
			text: handler.StartButton,
			setupMocks: func(t *mocks.MockTelegram, u *mocks.MockUsecase) {
				t.On("NewMessage", mock.Anything, mock.Anything).Times(1).Return(tgbotapi.MessageConfig{})
				t.On("Send", mock.Anything).Times(1).Return(tgbotapi.Message{}, nil)
			},
			exeptUser: &entity.User{
				CurrentStep: 0,
				IsCompleted: false,
				RemindStage: 0,
			},
		},
		{
			desc: "completed user gets completed message",
			user: &entity.User{
				IsCompleted: true,
			},
			exeptUser: &entity.User{
				IsCompleted: true,
			},
			text: "any",
			setupMocks: func(t *mocks.MockTelegram, u *mocks.MockUsecase) {
				t.On("NewMessage", mock.Anything, mock.Anything).Times(1).Return(tgbotapi.MessageConfig{})
				t.On("Send", mock.Anything).Times(1).Return(tgbotapi.Message{}, nil)
			},
		},
		{
			desc: "BackButton decrements step and resends question",
			user: &entity.User{

				CurrentStep: 2,
			},
			exeptUser: &entity.User{
				CurrentStep: 1,
			},
			text: handler.BackButton,
			setupMocks: func(t *mocks.MockTelegram, u *mocks.MockUsecase) {
				t.On("NewMessage", mock.Anything, mock.Anything).Times(1).Return(tgbotapi.MessageConfig{})
				t.On("Send", mock.Anything).Times(1).Return(tgbotapi.Message{}, nil)
			},
		},
		{
			desc: "user got finish survey",
			user: &entity.User{
				IsCompleted:  false,
				RemindStage:  handler.StartRemind,
				CurrentStep:  len(handler.Questions),
				EmailSentCnt: 0,
			},
			exeptUser: &entity.User{
				IsCompleted:  true,
				RemindStage:  handler.NotRemind,
				CurrentStep:  len(handler.Questions),
				EmailSentCnt: 1,
			},
			text: "any",
			setupMocks: func(t *mocks.MockTelegram, u *mocks.MockUsecase) {
				t.On("NewMessage", mock.Anything, mock.Anything).Times(1).Return(tgbotapi.MessageConfig{})
				t.On("Send", mock.Anything).Times(1).Return(tgbotapi.Message{}, nil)
				u.On("CreateEmail", mock.Anything).Times(1).Return(nil)
			},
		},
		{
			desc: "input is saved and advances step",
			user: &entity.User{
				CurrentStep: 0,
			},
			exeptUser: &entity.User{
				CurrentStep: 1,
			},
			text: "any",
			setupMocks: func(t *mocks.MockTelegram, u *mocks.MockUsecase) {
				t.On("NewMessage", mock.Anything, mock.Anything).Times(1).Return(tgbotapi.MessageConfig{})
				t.On("Send", mock.Anything).Times(1).Return(tgbotapi.Message{}, nil)
				u.On("SaveAnswer", mock.Anything).Times(1).Return(nil)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockTelegram := mocks.NewMockTelegram(t)
			mockUsecase := mocks.NewMockUsecase(t)
			if tC.setupMocks != nil {
				tC.setupMocks(mockTelegram, mockUsecase)
			}

			handler := handler.New(mockTelegram, zerolog.Nop(), mockUsecase, 1)

			handler.ProcessMessage(tC.user, tC.text)

			if diff := cmp.Diff(tC.exeptUser, tC.user, cmpopts.IgnoreFields(entity.User{}, "RemindAt")); diff != "" {
				t.Errorf("user mismatch (-want +got):\n%s", diff)
			}

			mockUsecase.AssertExpectations(t)
			mockTelegram.AssertExpectations(t)

		})
	}
}

func TestBotProcessContact(t *testing.T) {
	testCases := []struct {
		desc       string
		user       *entity.User
		text       string
		setupMocks func(telegram *mocks.MockTelegram, usecase *mocks.MockUsecase)
		exeptUser  *entity.User
	}{

		{
			desc: "completed user gets completed message",
			user: &entity.User{
				IsCompleted: true,
			},
			exeptUser: &entity.User{
				IsCompleted: true,
			},
			text: "any",
			setupMocks: func(t *mocks.MockTelegram, u *mocks.MockUsecase) {
				t.On("NewMessage", mock.Anything, mock.Anything).Times(1).Return(tgbotapi.MessageConfig{})
				t.On("Send", mock.Anything).Times(1).Return(tgbotapi.Message{}, nil)
			},
		},
		{
			desc: "input is saved and advances step",
			user: &entity.User{
				CurrentStep: 0,
			},
			exeptUser: &entity.User{
				CurrentStep: 1,
			},
			text: "any",
			setupMocks: func(t *mocks.MockTelegram, u *mocks.MockUsecase) {
				t.On("NewMessage", mock.Anything, mock.Anything).Times(1).Return(tgbotapi.MessageConfig{})
				t.On("Send", mock.Anything).Times(1).Return(tgbotapi.Message{}, nil)
				u.On("SaveAnswer", mock.Anything).Times(1).Return(nil)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockTelegram := mocks.NewMockTelegram(t)
			mockUsecase := mocks.NewMockUsecase(t)
			if tC.setupMocks != nil {
				tC.setupMocks(mockTelegram, mockUsecase)
			}

			handler := handler.New(mockTelegram, zerolog.Nop(), mockUsecase, 1)

			handler.ProcessContact(tC.user, tC.text)

			if diff := cmp.Diff(tC.exeptUser, tC.user, cmpopts.IgnoreFields(entity.User{}, "RemindAt")); diff != "" {
				t.Errorf("user mismatch (-want +got):\n%s", diff)
			}

			mockUsecase.AssertExpectations(t)
			mockTelegram.AssertExpectations(t)

		})
	}
}
