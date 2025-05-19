package app

import (
	"questionnaire-bot/internal/config"
	"questionnaire-bot/internal/repository"
	"questionnaire-bot/internal/telegram"
	"questionnaire-bot/internal/usecase"
	"questionnaire-bot/pkg/logger"
	"questionnaire-bot/pkg/postgres"
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
	//go monitoring.StartMetricsServer(os.Getenv("METRICS_SERVER_ADDR"))

	// usecase
	usecase := usecase.New(repository.New(pg))

	// workers

	bot, err := telegram.New(cfg.Bot.Token, l, usecase)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed initialized bot")
	}
	l.Info().Msg("Bot initialized")

	l.Info().Msg("Bot started")
	bot.Start()
}
