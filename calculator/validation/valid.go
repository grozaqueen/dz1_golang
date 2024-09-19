package validation

import (
	"unicode"

	"awesomeProject4/calc"
	"awesomeProject4/stackqueue"
)

// Дополнительная функция для проверки на оператор
func isOperator(symbol string) bool {
	return symbol == calc.OperatorAdd || symbol == calc.OperatorDivide || symbol == calc.OperatorMultiply || symbol == calc.OperatorSubtract
}

// Проверка на корректность числа (вещественное или целое)
func isNumber(str string) bool {
	// Проверка, содержит ли строка только цифры и одну точку
	dotCount := 0

	for i, r := range str {
		if !isValidCharacter(r) {
			return false
		}
		if string(r) == calc.SymbolDecimalPoint {
			dotCount++
			if dotCount > 1 || i == 0 || i == len(str)-1 {
				return false // Только одна точка допустима и не может быть в начале или в конце
			}
		}
	}
	return true
}

// Проверяет, является ли символ цифрой или точкой
func isValidCharacter(r rune) bool {
	return unicode.IsDigit(r) || string(r) == calc.SymbolDecimalPoint
}

func IsValidExpression(str string) bool {
	var (
		parenthesesStack stackqueue.Stack
		lastToken        rune
		currentNumber    string
	)

	for i, r := range str {
		symbol := string(r)

		switch {
		case symbol == calc.SymbolEmpty:
			continue // Пропускаем пробелы
		case unicode.IsDigit(r) || r == rune(calc.SymbolDecimalPoint[0]):
			if i > 0 && lastToken == rune(calc.SymbolRightParen[0]) {
				return false // Число после закрывающей скобки
			}
			currentNumber += symbol // Накопление числа
			lastToken = r
		case isOperator(symbol):
			if currentNumber != calc.SymbolEmpty {
				if !isNumber(currentNumber) {
					return false // Некорректное число (проверка на вещественное число)
				}
				currentNumber = calc.SymbolEmpty
			}
			if i == 0 || lastToken == rune(calc.SymbolLeftParen[0]) {
				if symbol != calc.SymbolUnaryMinus { // Унарный минус допустим
					return false // оператор в начале
				}
			}
			if isOperator(string(lastToken)) { // 2 оператора подряд
				return false
			}
			lastToken = r
		case symbol == calc.SymbolLeftParen:
			if lastToken == rune(calc.SymbolRightParen[0]) {
				return false // Оператор обязателен между скобками
			}
			parenthesesStack.Push(symbol)
			lastToken = r
		case symbol == calc.SymbolRightParen:
			if currentNumber != calc.SymbolEmpty {
				if !isNumber(currentNumber) {
					return false // Некорректное число перед закрывающей скобкой
				}
				currentNumber = calc.SymbolEmpty
			}
			if _, ok := parenthesesStack.Pop(); !ok {
				return false // Не хватает открывающей скобки
			}
			lastToken = r
		default:
			return false // Недопустимый символ
		}
	}

	if currentNumber != "" && !isNumber(currentNumber) {
		return false // Последнее число некорректно
	}

	// Проверка на сбалансированность скобок
	if len(parenthesesStack.Data) != 0 {
		return false
	}

	// Последний токен не должен быть оператором или открывающей скобкой
	if isOperator(string(lastToken)) || lastToken == rune(calc.SymbolLeftParen[0]) {
		return false
	}

	return true
}
