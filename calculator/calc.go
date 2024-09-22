package main

import (
	"fmt"
	"os"

	"awesomeProject4/calc"
	"awesomeProject4/preprocess"
	"awesomeProject4/validation"
)

const minArgs = 2

func main() {
	if len(os.Args) < minArgs {
		fmt.Println("Использование: go run calc.go \"выражение\"")
		return
	}

	str := os.Args[1]

	// Выполняем предобработку строки
	str = preprocess.PreprocessExpression(str)

	// Проверяем корректность выражения
	if !validation.IsValidExpression(str) {
		fmt.Println("Ошибка: Некорректное выражение")
		return
	}

	// Вызов функции calc для вычисления выражения
	result, err := calc.Calc(str)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println(result)
}
