package usecase

import (
	"fmt"
	"questionnaire-bot/internal/entity"
	"questionnaire-bot/internal/messages"
	"sort"
	"strings"
)

func (u *BotService) SaveAnswer(ans *entity.Answer) error {
	if err := u.repo.SaveAnswer(ans); err != nil {
		return err
	}
	return nil
}

func (u *BotService) GetComboMessage(id int64) (*messages.Combo, *messages.PersonalAdvice, error) {
	answers, err := u.repo.GetUserTechNames(id)
	if err != nil {
		return nil, nil, fmt.Errorf("BotService -> GetComboMessage -> repo.GetUserTechNames: %w", err)
	}

	if len(answers) != 5 {
		return nil, nil, fmt.Errorf("BotService -> GetComboMessage -> len(answers) != 5")
	}

	var personalAdvice *messages.PersonalAdvice
	advices := map[string]messages.PersonalAdvice{
		"q4a2": messages.Advice1,
		"//":   messages.Advice2, // todo ждем ответа клиента
		"q3a2": messages.Advice3,
	}
	for _, a := range answers {
		if advice, ok := advices[a.TechName]; ok {
			personalAdvice = &advice
		}
	}

	filtered := []string{
		answers[0].TechName,
		answers[1].TechName,
		answers[4].TechName,
	}
	sort.Strings(filtered)
	key := strings.Join(filtered, "|")

	combos := map[string]messages.Combo{
		"q1a1|q2a1|q5a2": messages.Combo1,
		"q1a3|q2a1|q5a1": messages.Combo2,
		"q1a2|q2a3|q5a2": messages.Combo3,
		"q1a1|q2a4|q5a1": messages.Combo4,
		"q1a2|q2a4|q5a1": messages.Combo4,
		"q1a1|q2a2|q5a3": messages.Combo5,
	}
	if msg, ok := combos[key]; ok {
		return &msg, personalAdvice, nil
	}
	return &messages.UniversalCombo, personalAdvice, nil
}
