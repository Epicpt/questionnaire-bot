package repository

import (
	"context"
	"questionnaire-bot/internal/entity"
)

const (
	queryGetChatIDsForNotify = `SELECT tg_id, chat_id, first_name, last_name, username, created_at, updated_at, remind_stage,remind_at,is_completed,current_step,max_step_reached FROM users WHERE remind_at < NOW() AND remind_stage < 2 AND is_completed = false`
)

func (r *BotRepo) GetUsersForNotify() ([]entity.User, error) {
	rows, err := r.Pool.Query(context.Background(), queryGetChatIDsForNotify)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.TgID, &user.ChatID, &user.FirstName, &user.LastName, &user.Username, &user.CreatedAt, &user.UpdatedAt,
			&user.RemindStage, &user.RemindAt, &user.IsCompleted, &user.CurrentStep, &user.MaxStepReached); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
