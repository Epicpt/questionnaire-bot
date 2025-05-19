package repository

import (
	"context"
	"questionnaire-bot/internal/entity"
)

const (
	querySaveEmail = `INSERT INTO emails (user_tg_id, body, status)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_tg_id) DO UPDATE SET body = $2, status = $3`
)

func (r *BotRepo) SaveEmail(e *entity.Email) error {
	_, err := r.Pool.Exec(context.Background(), querySaveEmail, e.UserTgID, e.Body, e.Status)
	if err != nil {
		return err
	}
	return nil
}
