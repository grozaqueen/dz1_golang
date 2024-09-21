package calc

import (
	"fmt"
	"strconv"
	"unicode"

	"awesomeProject4/stackqueue"
)

const (
	OperatorAdd      = "+"
	OperatorSubtract = "-"
	OperatorMultiply = "*"
	OperatorDivide   = "/"

	// Специальные символы
	SymbolDecimalPoint    = "."
	SymbolUnaryMinus      = "-"
	SymbolLeftParen       = "("
	SymbolRightParen      = ")"
	SymbolExpressionStart = "~"
	SymbolExpressionEnd   = "?"
	SymbolEmpty           = ""
	SymbolTab             = "\t"

	// Коды для матрицы операций
	MatrixPush    = 1
	MatrixReplace = 2
	MatrixPop     = 3
	MatrixEnqueue = 4
	MatrixError   = 100
	MatrixInvalid = -1
)

// Мапа для приоритета операций
var OperatorPrecedence = map[string][2]int{
	SymbolExpressionStart: {0, 0},
	OperatorAdd:           {1, 0},
	OperatorSubtract:      {1, 0},
	OperatorMultiply:      {2, 1},
	OperatorDivide:        {2, 1},
	SymbolLeftParen:       {3, 2},
	SymbolRightParen:      {3, 3},
	SymbolExpressionEnd:   {4, 4},
}

// Матрица операций
var operationMatrix = [4][5]int{
	{MatrixPush, MatrixPush, MatrixPush, MatrixError, MatrixInvalid},
	{MatrixReplace, MatrixPush, MatrixPush, MatrixEnqueue, MatrixEnqueue},
	{MatrixEnqueue, MatrixReplace, MatrixPush, MatrixEnqueue, MatrixEnqueue},
	{MatrixPush, MatrixPush, MatrixPush, MatrixPop, MatrixError},
}

// Функция получения индекса
func getIndex(symbol string) (int, int) {
	if precedence, exists := OperatorPrecedence[symbol]; exists {
		return precedence[0], precedence[1]
	}

	return MatrixInvalid, MatrixInvalid // Обработка некорректного символа
}

// Функция выполнения операций для чисел с плавающей точкой
func performOperation(op1, op2 float64, operator string) (float64, error) {
	switch operator {
	case OperatorAdd:
		return op1 + op2, nil
	case OperatorSubtract:
		return op1 - op2, nil
	case OperatorMultiply:
		return op1 * op2, nil
	case OperatorDivide:
		if op2 == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}
		return op1 / op2, nil
	default:
		return 0, fmt.Errorf("неизвестный оператор: %s", operator)
	}
}

func Calc(expression string) (float64, error) {
	// Начало и конец выражения
	str := SymbolExpressionStart + expression + SymbolExpressionEnd
	var (
		stackOperations stackqueue.Stack
		stackOperands   stackqueue.Stack
		queueOperands   stackqueue.Queue
		currentNumber   = SymbolEmpty
		isNegative      bool
	)

	stackOperations.Push(string(str[0]))

	for i, r := range str {
		if i == 0 {
			continue
		}

		if unicode.IsDigit(r) || string(r) == SymbolDecimalPoint {
			currentNumber += string(r)
		} else if string(r) == OperatorSubtract && (i == 1 || str[i-1] == SymbolLeftParen[0]) {
			isNegative = true
		} else {
			// Преобразуем и добавляем текущий операнд в очередь
			if currentNumber != SymbolEmpty {
				num := convertStringToNumber(currentNumber, isNegative)
				queueOperands.Enqueue(num)
				currentNumber = SymbolEmpty
				isNegative = false
			}

			// Обрабатываем операцию
			processed, err := processOperation(string(r), &stackOperations, &queueOperands, operationMatrix)
			if err != nil {
				return 0, err
			}

			// Продолжаем обработку, если операция требует повторной обработки (MatrixEnqueue)
			for processed {
				processed, err = processOperation(string(r), &stackOperations, &queueOperands, operationMatrix)
				if err != nil {
					return 0, err
				}
			}
		}
	}

	// Вынесено в функцию processRemainingOperations
	return processRemainingOperations(stackOperands, queueOperands)
}

// strToFloat64 преобразует строку в float64, если это возможно
func strToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// convertStringToNumber преобразует строку в число с учетом знака
func convertStringToNumber(currentNumber string, isNegative bool) float64 {
	if currentNumber != SymbolEmpty {
		if isNegative {
			currentNumber = SymbolUnaryMinus + currentNumber
		}

		num, _ := strToFloat64(currentNumber)

		return num
	}

	return 0
}

// processOperation обрабатывает символы операций с использованием матрицы приоритетов
func processOperation(symbol string,
	stackOperations *stackqueue.Stack,
	queueOperands *stackqueue.Queue,
	matrix [4][5]int) (bool, error) {

	value, _ := stackOperations.Peek()
	ind1, _ := getIndex(value.(string))
	_, ind2 := getIndex(symbol)

	switch matrix[ind1][ind2] {
	case MatrixPush:
		stackOperations.Push(symbol)
	case MatrixReplace:
		stackOperations.Pop()
		stackOperations.Push(symbol)
		queueOperands.Enqueue(value)
	case MatrixPop:
		stackOperations.Pop()
	case MatrixEnqueue:
		queueOperands.Enqueue(value)
		stackOperations.Pop()
		return true, nil // Повторяем обработку после удаления
	case MatrixError:
		return false, fmt.Errorf("ошибка в выражении")
	}

	return false, nil
}

// processRemainingOperations обрабатывает оставшиеся операнды и операции в очереди
func processRemainingOperations(stackOperands stackqueue.Stack, queueOperands stackqueue.Queue) (float64, error) {
	for len(queueOperands.Data) > 0 {
		element, ok := queueOperands.Dequeue()
		if !ok {
			break
		}

		switch v := element.(type) {
		case float64:
			stackOperands.Push(v)
		case string:
			op2, _ := stackOperands.Pop()
			op1, _ := stackOperands.Pop()

			op1Float := op1.(float64)
			op2Float := op2.(float64)

			// Выполняем операцию, включая деление на 0
			result, err := performOperation(op1Float, op2Float, v)
			if err != nil {

				return 0, err // Если возникает ошибка (например, деление на 0), возвращаем её
			}

			stackOperands.Push(result)
		}
	}

	result, _ := stackOperands.Pop()
	resultFloat, _ := result.(float64)

	return resultFloat, nil
}
