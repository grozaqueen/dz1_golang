package preprocess

import (
	"strings"

	"awesomeProject4/calc"
)

// preprocessExpression выполняет предобработку выражения
func PreprocessExpression(expression string) string {
	// Удаляем все пробелы и табуляции
	expression = strings.ReplaceAll(expression, calc.SymbolSpace, "")
	expression = strings.ReplaceAll(expression, calc.SymbolTab, "")

	// Предобработка: замена -( на (-1)*(
	if strings.HasPrefix(expression, calc.NegativeThreeUnaryMinus) {
		expression = calc.DenialBeforeNegativeThree + expression[1:] // Заменяем начало строки с -( на (-1)*(
	}

	return expression
}
