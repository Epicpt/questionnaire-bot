package app

import (
	"questionnaire-bot/internal/config"
	"questionnaire-bot/internal/handler"
	"questionnaire-bot/internal/repository"
	"questionnaire-bot/internal/usecase"
	"questionnaire-bot/internal/worker"
	"questionnaire-bot/pkg/logger"
	"questionnaire-bot/pkg/postgres"
	"questionnaire-bot/pkg/smtp"
	"questionnaire-bot/pkg/telegram"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	l.Info().Msg("Logger initialized")

	pg, err := postgres.New(cfg.PG.URL, cfg.PG.PoolMax)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer pg.Close()

	l.Info().Msg("PostgreSQL initialized")

	//metrics
	//go monitoring.StartMetricsServer(cfg)

	// usecase
	usecase := usecase.New(repository.New(pg), smtp.New(cfg.Smtp))

	botAPI, err := telegram.New(cfg.Bot.Token)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed initialized bot")
	}
	bot := handler.New(botAPI, l, usecase, cfg.Bot.AdminID)
	l.Info().Msg("Bot initialized")

	// worker
	worker := worker.New(usecase, l, cfg.Scheduler, bot)
	worker.Start()
	defer worker.Stop()
	l.Info().Msg("Worker started")

	l.Info().Msg("Bot started")
	bot.Start()
}
