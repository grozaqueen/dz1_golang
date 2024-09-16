package main

import (
	"awesomeProject4/calc"
	"awesomeProject4/stackqueue"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Дополнительная функция для проверки на оператор
func isOperator(symbol string) bool {
	switch symbol {
	case calc.OperatorAdd, calc.OperatorDivide, calc.OperatorMultiply, calc.OperatorSubtract:
		return true
	default:
		return false
	}
}

// Проверка на корректность числа (вещественное или целое)
func isNumber(str string) bool {
	// Число может содержать одну точку, перед которой и после которой должны быть цифры
	dotCount := 0
	for i, r := range str {
		if r == '.' {
			if i == 0 || i == len(str)-1 { // Точка не может быть в начале или в конце
				return false
			}
			dotCount++
			if dotCount > 1 { // Только одна точка допустима
				return false
			}
		} else if !unicode.IsDigit(r) { // Допустимы только цифры и точка
			return false
		}
	}
	return true
}

func isValidExpression(str string) bool {
	parenthesesStack := stackqueue.Stack{}
	var lastToken rune
	currentNumber := ""

	for i, r := range str {
		symbol := string(r)

		switch {
		case symbol == " ":
			continue // Пропускаем пробелы
		case unicode.IsDigit(r) || r == '.':
			if i > 0 && lastToken == ')' {
				return false // Число после закрывающей скобки
			}
			currentNumber += symbol // Накопление числа
			lastToken = r
		case isOperator(symbol):
			if currentNumber != "" {
				if !isNumber(currentNumber) {
					return false // Некорректное число (проверка на вещественное число)
				}
				currentNumber = ""
			}
			if i == 0 || lastToken == '(' {
				if symbol != "-" { // Унарный минус допустим
					return false // оператор в начале
				}
			}
			if isOperator(string(lastToken)) { // 2 оператора подряд
				return false
			}
			lastToken = r
		case symbol == "(":
			if lastToken == ')' {
				return false // Оператор обязателен между скобками
			}
			parenthesesStack.Push(symbol)
			lastToken = r
		case symbol == ")":
			if currentNumber != "" {
				if !isNumber(currentNumber) {
					return false // Некорректное число перед закрывающей скобкой
				}
				currentNumber = ""
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
	if len(parenthesesStack) != 0 {
		return false
	}

	// Последний токен не должен быть оператором или открывающей скобкой
	if isOperator(string(lastToken)) || lastToken == '(' {
		return false
	}

	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run calc.go \"выражение\"")
		return
	}

	str := os.Args[1]

	// Удаляем все пробелы и табуляции
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\t", "")

	if !isValidExpression(str) {
		fmt.Println("Ошибка: Некорректное выражение")
		return
	}

	result, err := calc.Calc(str) // Вызов функции calc
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println(result)
}
