package config

import (
	c "questionnaire-bot/pkg/configuration"
	"time"
)

type Config struct {
	Log       c.Log
	PG        c.PG
	Bot       Bot
	Smtp      SMTP
	Scheduler Scheduler
}

type Bot struct {
	Token         string `env-required:"true" env:"TELEGRAM_BOT_TOKEN"`
	EmployeesData EmployeesData
}

type EmployeesData struct {
	AdminID    int64    `env-required:"true" env:"TELEGRAM_ADMIN_ID"`
	ManagerIDs []string `env-required:"true" env:"TELEGRAM_MANAGER_IDS"`
}

type SMTP struct {
	Host     string `env-required:"true" env:"SMTP_HOST"`
	Port     int    `env-required:"true" env:"SMTP_PORT"`
	Username string `env-required:"true" env:"SMTP_USERNAME"`
	Password string `env-required:"true" env:"SMTP_PASSWORD"`
	Client   string `env-required:"true" env:"SMTP_CLIENT"`
	Sender   string `env-required:"true" env:"SMTP_SENDER"`
}

type Scheduler struct {
	EmailSend    time.Duration `env-required:"true" env:"SCHEDULER_EMAIL_SEND_PERIOD"`
	EmailErrSend time.Duration `env-required:"true" env:"SCHEDULER_RETRY_ERROR_SEND_PERIOD"`
	NotifySend   time.Duration `env-required:"true" env:"SCHEDULER_NOTIFICATION_SEND_PERIOD"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := c.Load(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
