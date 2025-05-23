package tests

import (
	"strings"
	"testing"

	"questionnaire-bot/internal/handler"

	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	testCases := []struct {
		desc        string
		input       string
		expectedErr error
	}{
		{
			desc:        "Valid name",
			input:       "Иван Петров",
			expectedErr: nil,
		},
		{
			desc:        "Empty name",
			input:       "",
			expectedErr: handler.EmptyNameErr,
		},
		{
			desc:        "Name with numbers",
			input:       "Иван123",
			expectedErr: handler.RegexpNameErr,
		},
		{
			desc:        "Name with special chars",
			input:       "Иван@Петров",
			expectedErr: handler.RegexpNameErr,
		},
		{
			desc:        "Name with extra special chars",
			input:       "Иван-Петров ():SQL?/```~",
			expectedErr: handler.RegexpNameErr,
		},
		{
			desc:        "long name",
			input:       strings.Repeat("A", 101),
			expectedErr: handler.LongNameErr,
		},
		{
			desc:        "Unicode letters",
			input:       "Élève Zürich",
			expectedErr: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := handler.ValidateName(tC.input)
			if tC.expectedErr != nil {
				assert.EqualError(t, err, tC.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestValidateArea(t *testing.T) {
	testCases := []struct {
		desc        string
		input       string
		expectedErr error
	}{
		{
			desc:        "Valid area",
			input:       "100",
			expectedErr: nil,
		},
		{
			desc:        "Empty area",
			input:       "",
			expectedErr: handler.EmptyAreaErr,
		},
		{
			desc:        "Not an area",
			input:       "abc",
			expectedErr: handler.ConvAreaErr,
		},
		{
			desc:        "Negative value",
			input:       "-123",
			expectedErr: handler.NegativeAreaErr,
		},
		{
			desc:        "Below minimum (39)",
			input:       "39",
			expectedErr: handler.Less40AreaErr,
		},
		{
			desc:        "Minimum boundary (40)",
			input:       "40",
			expectedErr: nil,
		},
		{
			desc:        "Normal value",
			input:       "85",
			expectedErr: nil,
		},
		{
			desc:        "Very large number",
			input:       "100001",
			expectedErr: handler.OverAreaErr,
		},
		{
			desc:        "Maximum boundary",
			input:       "100000",
			expectedErr: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := handler.ValidateArea(tC.input)
			if tC.expectedErr != nil {
				assert.EqualError(t, err, tC.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestValidateCity(t *testing.T) {
	testCases := []struct {
		desc        string
		input       string
		expectedErr error
	}{
		{
			desc:        "Empty input",
			input:       "   ",
			expectedErr: handler.CityEmptyErr,
		},
		{
			desc:        "Empty string",
			input:       "",
			expectedErr: handler.CityEmptyErr,
		},

		{
			desc:        "Too long city name",
			input:       strings.Repeat("A", 101),
			expectedErr: handler.LongCityErr,
		},
		{
			desc:        "Max allowed length (100 chars)",
			input:       strings.Repeat("A", 100),
			expectedErr: nil,
		},

		{
			desc:        "Simple city name",
			input:       "Москва",
			expectedErr: nil,
		},
		{
			desc:        "City with multiple spaces (should be normalized)",
			input:       "  Санкт   Петербург  ",
			expectedErr: nil,
		},
		{
			desc:        "Unicode city name",
			input:       "Çeñàbà",
			expectedErr: nil,
		},
		{
			desc:        "City with numbers",
			input:       "Москва123",
			expectedErr: handler.RegexpCityErr,
		},
		{
			desc:        "City with special chars",
			input:       "Москва@",
			expectedErr: handler.RegexpCityErr,
		},
		{
			desc:        "City with brackets",
			input:       "Москва(центральная)",
			expectedErr: handler.RegexpCityErr,
		},
		{
			desc:        "City with dot",
			input:       "Санкт.Петербург",
			expectedErr: handler.RegexpCityErr,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := handler.ValidateCity(tC.input)
			if tC.expectedErr != nil {
				assert.EqualError(t, err, tC.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestValidatePhone(t *testing.T) {
	testCases := []struct {
		desc        string
		input       string
		expectedErr error
	}{
		{
			desc:        "Empty input",
			input:       "   ",
			expectedErr: handler.EmptyPhoneErr,
		},
		{
			desc:        "Empty string",
			input:       "",
			expectedErr: handler.EmptyPhoneErr,
		},

		{
			desc:        "Valid 8XXXXXXXXXX",
			input:       "89161234567",
			expectedErr: nil,
		},
		{
			desc:        "Valid +7XXXXXXXXXX",
			input:       "+79161234567",
			expectedErr: nil,
		},
		{
			desc:        "Valid 7XXXXXXXXXX",
			input:       "79161234567",
			expectedErr: nil,
		},
		{
			desc:        "With spaces",
			input:       "+7 916 123 45 67",
			expectedErr: nil,
		},
		{
			desc:        "With hyphens",
			input:       "8(916)123-45-67",
			expectedErr: nil,
		},
		{
			desc:        "With parentheses",
			input:       "7(916)1234567",
			expectedErr: nil,
		},
		{
			desc:        "Mixed formatting",
			input:       "+7 (916) 123-45-67",
			expectedErr: nil,
		},
		{
			desc:        "Too short",
			input:       "8916123456",
			expectedErr: handler.FormatPhoneErr,
		},
		{
			desc:        "Too long",
			input:       "891612345678",
			expectedErr: handler.FormatPhoneErr,
		},
		{
			desc:        "Invalid prefix",
			input:       "+59161234567",
			expectedErr: handler.FormatPhoneErr,
		},
		{
			desc:        "Letters in number",
			input:       "8abc1234567",
			expectedErr: handler.FormatPhoneErr,
		},
		{
			desc:        "Special characters",
			input:       "8@916!234#67",
			expectedErr: handler.FormatPhoneErr,
		},
		{
			desc:        "International format",
			input:       "+442012345678",
			expectedErr: handler.FormatPhoneErr,
		},
		{
			desc:        "Shortcode",
			input:       "123",
			expectedErr: handler.FormatPhoneErr,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := handler.ValidatePhone(tC.input)
			if tC.expectedErr != nil {
				assert.EqualError(t, err, tC.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
