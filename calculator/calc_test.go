package main

import (
	"testing"

	"awesomeProject4/calc"
	"awesomeProject4/validation"
	"github.com/stretchr/testify/require"
)

func TestIsValidExpression_Success_Advanced(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected bool
	}{
		{"(2+3)*(5-2)", true},
		{"((2+3)*(4-1))/3", true},
		{"-1+(-2*(-3))", true},
		{"3.5+2.1", true}, // Проверка на дробные числа
		{"(1+2)*(3-(4/2))", true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			result := validation.IsValidExpression(test.input)

			require.Equal(t, test.expected, result, "Для %s, ожидалось %v, но имеем %v", test.input, test.expected, result)
		})
	}
}

func TestIsValidExpression_Failure_Advanced(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected bool
	}{
		{"1++2", false},      // Двойной оператор
		{"(2+3)*(4-", false}, // Недостающая закрывающая скобка
		{"2*(3-4))", false},  // Лишняя закрывающая скобка
		{"3.5.1 + 2", false}, // Некорректное дробное число
		{"()()", false},      // Скобки без операндов и операторов
		{"+2*3", false},      // Оператор в начале без унарного минуса
		{"(2+3)()", false},   // Оператор обязателен между скобками
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			result := validation.IsValidExpression(test.input)

			require.Equal(t, test.expected, result, "Для %s, ожидалось %v, но имеем %v", test.input, test.expected, result)
		})
	}
}

func TestCalculator_Success_Advanced(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected float64
	}{
		{"2+2", 4},
		{"(2+3)*(5-2)", 15},
		{"10/2+3*2", 11},
		{"-10+5", -5},
		{"3.5+2.1", 5.6}, // Проверка работы с дробями
		{"(1+2)*(3-(4/2))", 3},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			result, err := calc.Calc(test.input)

			require.NoError(t, err, "Для %s, неожиданная ошибка: %v", test.input, err)
			require.InEpsilon(t, test.expected, result, 0.0001, "Для %s, ожидалось %f, но имеем %f", test.input, test.expected, result)
		})
	}
}

func TestCalculator_Failure_Advanced(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
	}{
		{"10/(5-5)"}, // Деление на ноль
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			_, err := calc.Calc(test.input)

			require.Error(t, err, "Для %s, ожидалась ошибка, но ничего не вышло", test.input)
		})
	}
}
