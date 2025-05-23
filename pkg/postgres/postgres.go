package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// const uniqueViolationCode = "23505"

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(dsn string, maxConn int) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга строки подключения: %w", err)
	}

	config.MaxConns = int32(maxConn)

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пула соединений: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("база данных недоступна: %w", err)
	}

	return &Postgres{Pool: pool}, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
		log.Info().Msg("Подключение к PostgreSQL закрыто")
	}
}

func IsNotFoundErr(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

// func IsUniqueViolation(err error) bool {
// 	var pgErr *pgconn.PgError
// 	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
// }

// func IsTxClosed(err error) bool {
// 	return errors.Is(err, pgx.ErrTxClosed)
// }
