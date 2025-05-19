package telegram

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	validateName = func(input string) error {
		input = strings.TrimSpace(input)
		input = strings.Join(strings.Fields(input), " ")

		if len(input) == 0 {
			return errors.New("Имя не должно быть пустым")
		}
		if len(input) > 100 {
			return errors.New("Имя слишком длинное")
		}

		ok, err := regexp.MatchString(`^[\p{L} ]+$`, input)
		if err != nil {
			return errors.New("Ошибка проверки имени")
		}
		if !ok {
			return errors.New("Имя должно содержать только буквы и пробелы")
		}
		return nil
	}

	validateArea = func(input string) error {
		input = strings.TrimSpace(input)

		if len(input) == 0 {
			return errors.New("Площадь не должна быть пустой")
		}

		n, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("Площадь должна быть числом")
		}

		if n <= 0 {
			return errors.New("Площадь должна быть положительным числом")
		}
		if n < 40 {
			return errors.New(`Спасибо за уточнение!
Минимальная площадь для разработки дизайн-проекта в нашей студии - 40 кв.м. Возможно у вас есть дополнительная площадь для которой Вы хотели бы разработать проект?`)
		}

		if n > 100000 {
			return errors.New("Площадь должна быть меньше")
		}

		return nil
	}

	validateObject = func(s string) error {
		for _, opt := range objectOpt {
			if s == opt {
				return nil
			}
		}
		return errors.New("Выберите один из предложенных вариантов")
	}

	validateRenovation = func(s string) error {
		for _, opt := range renovationOpt {
			if s == opt {
				return nil
			}
		}
		return errors.New("Выберите один из предложенных вариантов")
	}

	validateEquipment = func(s string) error {
		for _, opt := range equipmentOpt {
			if s == opt {
				return nil
			}
		}
		return errors.New("Выберите один из предложенных вариантов")
	}

	validateConnection = func(s string) error {
		for _, opt := range connectionOpt {
			if s == opt {
				return nil
			}
		}
		return errors.New("Выберите один из предложенных вариантов")
	}

	validateCity = func(input string) error {
		input = strings.TrimSpace(input)
		input = strings.Join(strings.Fields(input), " ")

		if len(input) == 0 {
			return errors.New("Город не должно быть пустым")
		}
		if len(input) > 100 {
			return errors.New("Название города слишком длинное")
		}

		ok, err := regexp.MatchString(`^[\p{L}’'\- ]+$`, input)
		if err != nil {
			return errors.New("Ошибка проверки города")
		}
		if !ok {
			return errors.New("Название города может содержать только буквы, пробелы, апострофы и дефис")
		}
		return nil
	}

	validatePhone = func(input string) error {
		input = strings.TrimSpace(input)
		input = strings.NewReplacer(" ", "", "(", "", ")", "", "-", "").Replace(input)

		if len(input) == 0 {
			return errors.New("Телефон не должен быть пустым")
		}

		re := regexp.MustCompile(`^(?:\+7|8|7)\d{10}$`)
		if !re.MatchString(input) {
			return errors.New("Введите корректный номер телефона в формате +7XXXXXXXXXX или 8XXXXXXXXXX")
		}

		return nil
	}
)
