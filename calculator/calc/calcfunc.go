package calc

import (
	"awesomeProject4/stackqueue"
	"fmt"
	"strconv"
	"unicode"
)

const (
	OperatorAdd      = "+"
	OperatorSubtract = "-"
	OperatorMultiply = "*"
	OperatorDivide   = "/"
)

// Мапа для приоритета операций
var OperatorPrecedence = map[string][2]int{
	"~":              {0, 0},
	OperatorAdd:      {1, 0},
	OperatorSubtract: {1, 0},
	OperatorMultiply: {2, 1},
	OperatorDivide:   {2, 1},
	"(":              {3, 2},
	")":              {3, 3},
	"?":              {4, 4},
}

// Функция получения индекса
func getIndex(symbol string) (int, int) {
	if precedence, exists := OperatorPrecedence[symbol]; exists {
		return precedence[0], precedence[1]
	}
	return -1, -1 // Обработка некорректного символа
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
	matrix := [4][5]int{
		{1, 1, 1, 100, -1},
		{2, 1, 1, 4, 4},
		{4, 2, 1, 4, 4},
		{1, 1, 1, 3, 100},
	}

	begin := "~"
	end := "?"
	str := begin + expression + end

	var stackOperations stackqueue.Stack
	var stackOperands stackqueue.Stack
	var queueOperands stackqueue.Queue
	currentNumber := ""
	var isNegative bool
	stackOperations.Push(string(str[0]))
	for i, r := range str {
		if i == 0 {
			continue
		}
	repeat:
		if unicode.IsDigit(r) || r == '.' {
			currentNumber += string(r)
		} else if r == '-' && (i == 1 || str[i-1] == '(') {
			isNegative = true
		} else {
			if currentNumber != "" {
				if isNegative {
					currentNumber = "-" + currentNumber
					isNegative = false
				}
				num, _ := strconv.ParseFloat(currentNumber, 64)
				queueOperands.Enqueue(num)
				currentNumber = ""
			}

			symbol := string(r)
			value, _ := stackOperations.Peek()
			ind1, _ := getIndex(value.(string))
			_, ind2 := getIndex(symbol)

			switch matrix[ind1][ind2] {
			case 1:
				stackOperations.Push(symbol)
			case 2:
				stackOperations.Pop()
				stackOperations.Push(symbol)
				queueOperands.Enqueue(value)
			case 3:
				stackOperations.Pop()
			case 4:
				queueOperands.Enqueue(value)
				stackOperations.Pop()
				goto repeat
			case 100:
				return 0, fmt.Errorf("ошибка в выражении")
			}
		}
	}

	for len(queueOperands) > 0 {
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

			result, err := performOperation(op1Float, op2Float, v)
			if err != nil {
				return 0, err // Возврат ошибки
			}
			stackOperands.Push(result)
		}
	}

	result, _ := stackOperands.Pop()
	resultFloat, _ := result.(float64)
	return resultFloat, nil
}
