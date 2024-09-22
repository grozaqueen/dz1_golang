package main

import (
	"fmt"
	"os"
	"strings"

	"awesomeProject4/calc"
	"awesomeProject4/validation"
)

const minArgs = 2

func main() {
	if len(os.Args) < minArgs {
		fmt.Println("Использование: go run calc.go \"выражение\"")
		return
	}

	str := os.Args[1]

	// Удаляем все пробелы и табуляции
	str = strings.ReplaceAll(str, calc.SymbolSpace, "")
	str = strings.ReplaceAll(str, calc.SymbolTab, "")
	if !validation.IsValidExpression(str) {
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
