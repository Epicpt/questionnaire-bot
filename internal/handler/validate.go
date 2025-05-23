package handler

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	EmptyNameErr  = errors.New("Имя не должно быть пустым")
	LongNameErr   = errors.New("Имя слишком длинное")
	RegexpNameErr = errors.New("Имя должно содержать только буквы и пробелы")

	EmptyAreaErr    = errors.New("Площадь не должна быть пустой")
	NegativeAreaErr = errors.New("Площадь должна быть положительным числом")
	ConvAreaErr     = errors.New("Площадь должна быть числом")
	Less40AreaErr   = errors.New(`Спасибо за уточнение!
Минимальная площадь для разработки дизайн-проекта в нашей студии - 40 кв.м. Возможно у вас есть дополнительная площадь для которой Вы хотели бы разработать проект?`)
	OverAreaErr = errors.New("Площадь должна быть меньше")

	ChooseErr = errors.New("Выберите один из предложенных вариантов")

	CityEmptyErr  = errors.New("Город не должен быть пустым")
	LongCityErr   = errors.New("Название города слишком длинное")
	RegexpCityErr = errors.New("Название города может содержать только буквы, пробелы, апострофы и дефис")

	EmptyPhoneErr  = errors.New("Телефон не должен быть пустым")
	FormatPhoneErr = errors.New("Введите корректный номер телефона в формате +7XXXXXXXXXX или 8XXXXXXXXXX")
)

var (
	ValidateName = func(input string) error {
		input = strings.TrimSpace(input)
		input = strings.Join(strings.Fields(input), " ")

		if len(input) == 0 {
			return EmptyNameErr
		}
		if len(input) > 100 {
			return LongNameErr
		}

		ok, err := regexp.MatchString(`^[\p{L} ]+$`, input)
		if err != nil {
			return errors.New("Ошибка проверки имени")
		}
		if !ok {
			return RegexpNameErr
		}
		return nil
	}

	ValidateArea = func(input string) error {
		input = strings.TrimSpace(input)

		if len(input) == 0 {
			return EmptyAreaErr
		}

		n, err := strconv.Atoi(input)
		if err != nil {
			return ConvAreaErr
		}

		if n < 0 {
			return NegativeAreaErr
		}

		if n < 40 {
			return Less40AreaErr
		}

		if n > 100000 {
			return OverAreaErr
		}

		return nil
	}

	ValidateCity = func(input string) error {
		input = strings.TrimSpace(input)
		input = strings.Join(strings.Fields(input), " ")

		if len(input) == 0 {
			return CityEmptyErr
		}
		if len(input) > 100 {
			return LongCityErr
		}

		ok, err := regexp.MatchString(`^[\p{L}’'\- ]+$`, input)
		if err != nil {
			return errors.New("Ошибка проверки города")
		}
		if !ok {
			return RegexpCityErr
		}
		return nil
	}

	ValidatePhone = func(input string) error {
		input = strings.TrimSpace(input)
		input = strings.NewReplacer(" ", "", "(", "", ")", "", "-", "").Replace(input)

		if len(input) == 0 {
			return EmptyPhoneErr
		}

		re := regexp.MustCompile(`^(?:\+7|8|7)\d{10}$`)
		if !re.MatchString(input) {
			return FormatPhoneErr
		}

		return nil
	}
)

func ValidateChoose(s string, options []string) error {
	for _, opt := range options {
		if s == opt {
			return nil
		}
	}
	return ChooseErr
}
