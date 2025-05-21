package repository

import (
	"context"
	"questionnaire-bot/internal/entity"
)

const (
	querySaveEmail = `INSERT INTO emails (user_tg_id, body, status)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_tg_id) DO UPDATE SET body = $2, status = $3, updated_at = NOW() `

	queryGetEmailsByStatus = `SELECT id, user_tg_id, body, status
	FROM emails
	WHERE status = $1
	`

	queryUpdateStatus = `UPDATE emails SET status = $1, updated_at = NOW() WHERE id = $2`
)

func (r *BotRepo) SaveEmail(e *entity.Email) error {
	_, err := r.Pool.Exec(context.Background(), querySaveEmail, e.UserTgID, e.Body, e.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *BotRepo) GetEmailsByStatus(status string) ([]entity.Email, error) {
	rows, err := r.Pool.Query(context.Background(), queryGetEmailsByStatus, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []entity.Email
	for rows.Next() {
		var email entity.Email
		if err := rows.Scan(&email.ID, &email.UserTgID, &email.Body, &email.Status); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return emails, nil
}

func (r *BotRepo) UpdateEmailStatus(email *entity.Email, status string) error {
	_, err := r.Pool.Exec(context.Background(), queryUpdateStatus, status, email.ID)
	if err != nil {
		return err
	}
	return nil
}
