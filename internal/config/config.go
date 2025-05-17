package config

import (
	c "questionnaire-bot/pkg/configuration"
)

type Config struct {
	Log c.Log
	PG  c.PG
	Bot Bot
}

type Bot struct {
	Token string `env-required:"true" env:"TELEGRAM_BOT_TOKEN"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := c.Load(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
