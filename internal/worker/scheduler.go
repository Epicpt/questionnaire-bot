package worker

import (
	"fmt"
	"questionnaire-bot/internal/config"
	"questionnaire-bot/internal/constantses"
	"questionnaire-bot/internal/handler"
	"questionnaire-bot/internal/usecase"
	"time"

	"github.com/rs/zerolog"
)

const (
	sent    = "sent"
	failed  = "failed"
	pending = "pending"

	remindAfterDay = 1
	day            = 24 * time.Hour
)

type Scheduler struct {
	u           usecase.Usecase
	l           zerolog.Logger
	tickers     tickers
	stopChannel chan struct{}
	notifier    handler.Handler
	empCfg      config.EmployeesData
}

type tickers struct {
	email       *time.Ticker
	failedEmail *time.Ticker
	notify      *time.Ticker
}

func New(u usecase.Usecase, l zerolog.Logger, cfg config.Scheduler, notifier handler.Handler, empCfg config.EmployeesData) *Scheduler {
	return &Scheduler{
		u: u,
		l: l,
		tickers: tickers{
			email:       time.NewTicker(cfg.EmailSend),
			failedEmail: time.NewTicker(cfg.EmailErrSend),
			notify:      time.NewTicker(cfg.NotifySend)},
		stopChannel: make(chan struct{}),
		notifier:    notifier,
		empCfg:      empCfg,
	}
}

func (s *Scheduler) Start() {
	go func() {

		defer func() {
			if r := recover(); r != nil {
				s.l.Warn().Interface("panic", r).Msg("Паника в планировщике")
				// TODO: add metric here "panic in scheduler"
				time.Sleep(5 * time.Second)
				go s.Start()
			}
		}()

		for {
			select {
			//case <-s.tickers.email.C:
			//	s.processEmails()
			//case <-s.tickers.failedEmail.C:
			//	s.processFailedEmails()
			case <-s.tickers.notify.C:
				s.processNotify()

			case <-s.stopChannel:
				return
			}
		}

	}()
}

func (s *Scheduler) Stop() {
	if s.tickers.email != nil {
		s.tickers.email.Stop()
	}
	if s.tickers.failedEmail != nil {
		s.tickers.failedEmail.Stop()
	}
	if s.tickers.notify != nil {
		s.tickers.notify.Stop()
	}
	close(s.stopChannel)
	s.l.Info().Msg("Scheduler stopped")
}

func (s *Scheduler) processEmails() {
	emails, err := s.u.GetEmailsByStatus(pending)
	if err != nil {
		s.l.Error().Err(err).Msg("failed to get pending emails")
	}

	for _, email := range emails {
		status := sent
		if err := s.u.SendEmail(&email); err != nil {
			s.l.Error().Err(err).Int("email_id", email.ID).Msg("failed to send email")
			status = failed
		} else {
			s.l.Info().Int("email_id", email.ID).Msg("email sent")
		}

		if err := s.u.UpdateEmailStatus(&email, status); err != nil {
			s.l.Error().Err(err).Msg("failed to update email status")
		}
	}
}

func (s *Scheduler) processFailedEmails() {
	emails, err := s.u.GetEmailsByStatus(failed)
	if err != nil {
		s.l.Error().Err(err).Msg("failed to get failed emails")
	}

	for _, email := range emails {
		status := sent
		if err := s.u.SendEmail(&email); err != nil {
			s.l.Error().Err(err).Int("email_id", email.ID).Msg("failed to send email")
			status = failed
			msg := fmt.Sprintf("❗️Повторная отправка письма не удалась❗️\nВозможно SMTP сервер упал.\nID: %d", email.ID)
			s.notifier.SendTo(s.empCfg.AdminID, msg)
		} else {
			s.l.Info().Int("email_id", email.ID).Msg("email sent")
		}

		if err := s.u.UpdateEmailStatus(&email, status); err != nil {
			s.l.Error().Err(err).Msg("failed to update email status")
		}
	}
}

func (s *Scheduler) processNotify() {
	users, err := s.u.GetUsersForNotify()
	if err != nil {
		s.l.Error().Err(err).Msg("failed to get users for notify")
	}

	for _, user := range users {
		q := handler.NotifyQuestion

		s.notifier.Send(user.ChatID, q.Text, handler.KeyboardFromOptions(q, handler.ShowBackButton(user.CurrentStep)))
		user.RemindStage = constantses.NotRemind

		if err = s.u.SaveUser(&user); err != nil {
			s.l.Error().Err(err).Int64("id", user.TgID).Msg("Ошибка при сохранении пользователя")
		}
	}
	if len(users) > 0 {
		s.l.Info().Msgf("Notify sent to %d users", len(users))
	}
}
