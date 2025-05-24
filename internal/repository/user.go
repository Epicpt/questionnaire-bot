package repository

import (
	"context"
	"fmt"
	"questionnaire-bot/internal/entity"
	"questionnaire-bot/pkg/postgres"
)

var _ Repository = (*BotRepo)(nil)

const (
	querySaveUser = `
		INSERT INTO users (tg_id, chat_id, first_name, last_name, username, remind_stage,remind_at,is_completed,current_step,max_step_reached,email_sent_count) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		ON CONFLICT (tg_id) DO UPDATE SET updated_at = NOW(), remind_stage = $6, remind_at = $7, is_completed = $8, current_step = $9, max_step_reached = $10, email_sent_count = $11`

	queryGetUser = `
		SELECT tg_id, chat_id, first_name, last_name, username, created_at, updated_at, remind_stage,remind_at,is_completed,current_step,max_step_reached,email_sent_count
		FROM users
		WHERE tg_id = $1`
)

func (r *BotRepo) SaveUser(u *entity.User) error {
	_, err := r.Pool.Exec(context.Background(), querySaveUser,
		u.TgID, u.ChatID, u.FirstName, u.LastName, u.Username, u.RemindStage, u.RemindAt, u.IsCompleted, u.CurrentStep, u.MaxStepReached, u.EmailSentCnt)
	if err != nil {
		return err
	}

	return nil
}

func (r *BotRepo) GetUser(userID int64) (*entity.User, error) {
	var user entity.User

	err := r.Pool.QueryRow(context.Background(), queryGetUser, userID).Scan(&user.TgID, &user.ChatID, &user.FirstName,
		&user.LastName, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.RemindStage, &user.RemindAt, &user.IsCompleted, &user.CurrentStep, &user.MaxStepReached, &user.EmailSentCnt)

	if err != nil {
		if postgres.IsNotFoundErr(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("ошибка получения пользователя из БД: %w", err)
	}

	return &user, nil
}
