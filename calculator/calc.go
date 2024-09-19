package main

import (
	"fmt"
	"os"
	"strings"

	"awesomeProject4/calc"
	"awesomeProject4/validation"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run calc.go \"выражение\"")
		return
	}

	str := os.Args[1]

	// Удаляем все пробелы и табуляции
	str = strings.ReplaceAll(str, calc.SymbolEmpty, "")
	str = strings.ReplaceAll(str, "\t", "")

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
