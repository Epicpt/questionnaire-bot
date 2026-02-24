package repository

import (
	"context"
	"questionnaire-bot/internal/entity"
)

const (
	querySaveAnswer = `INSERT INTO answers (user_tg_id, question_key, step, answer, short, tech_name)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_tg_id, question_key, step) DO UPDATE SET answer = $4, tech_name = $6`
	queryGetAnswersUser = `SELECT step, question_key, answer, short
FROM answers
WHERE user_tg_id = $1
ORDER BY step ASC`
	queryGetUserTechNames = `SELECT tech_name FROM answers WHERE user_tg_id = $1 ORDER BY step ASC`
)

func (r *BotRepo) SaveAnswer(ans *entity.Answer) error {
	_, err := r.Pool.Exec(context.Background(), querySaveAnswer, ans.UserTgID, ans.QuestionKey, ans.Step, ans.UserAnswer, ans.Short, ans.TechName)
	if err != nil {
		return err
	}
	return nil
}

func (r *BotRepo) GetAnswersUser(tgID int64) ([]entity.Answer, error) {
	rows, err := r.Pool.Query(context.Background(), queryGetAnswersUser, tgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []entity.Answer
	for rows.Next() {
		var ans entity.Answer
		if err := rows.Scan(&ans.Step, &ans.QuestionKey, &ans.UserAnswer, &ans.Short); err != nil {
			return nil, err
		}
		answers = append(answers, ans)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}

func (r *BotRepo) GetUserTechNames(tgID int64) ([]entity.Answer, error) {
	rows, err := r.Pool.Query(context.Background(), queryGetUserTechNames, tgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []entity.Answer
	for rows.Next() {
		var ans entity.Answer
		if err := rows.Scan(&ans.TechName); err != nil {
			return nil, err
		}
		answers = append(answers, ans)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}
