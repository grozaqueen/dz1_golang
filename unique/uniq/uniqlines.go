package uniq

import (
	"fmt"
	"strings"
)

type Options struct {
	Mode       string // режим работы (default, count, duplicate, unique)
	IgnoreCase bool   // игнорировать регистр букв
	NumFields  int    // количество игнорируемых полей
	NumChars   int    // количество игнорируемых символов
}

// Функция для разделения строки на поля
func splitIntoFields(s string) []string {
	var fields []string
	field := ""
	inField := false

	for _, r := range s {
		if r == ' ' || r == '\t' { // Пробел или табуляция считаются разделителями
			if inField { // Если мы были внутри поля, то конец поля
				fields = append(fields, field)
				field = ""
				inField = false
			}
		} else {
			inField = true
			field += string(r) // Добавляем символ к полю
		}
	}

	// Добавляем последнее поле, если строка не закончилась разделителем
	if inField {
		fields = append(fields, field)
	}

	return fields
}

// функция сравнения пары строк с учетом введенных параметров
func CompareStrings(str1, str2 string, ignoreCase bool, numFields, numChars int) bool {
	// Игнорируем регистр, если установлен флаг
	if ignoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	// Если нужно игнорировать numFields полей
	if numFields > 0 {
		fields1 := splitIntoFields(str1)
		fields2 := splitIntoFields(str2)

		// Обрезаем поля для str1
		if len(fields1) > numFields {
			str1 = ""
			for i := numFields; i < len(fields1); i++ {
				if i > numFields {
					str1 += " " // Вставляем пробелы между полями
				}
				str1 += fields1[i]
			}
		}

		// Аналогично для str2
		if len(fields2) > numFields {
			str2 = ""
			for i := numFields; i < len(fields2); i++ {
				if i > numFields {
					str2 += " "
				}
				str2 += fields2[i]
			}
		}
	}

	// Игнорируем numChars символов путем среза
	if numChars > 0 {
		if len(str1) >= numChars {
			str1 = str1[numChars:]
		}
		if len(str2) >= numChars {
			str2 = str2[numChars:]
		}
	}

	// Сравниваем строки
	return str1 == str2
}

// Функция для уникализации строк
func ProcessStrings(lines []string, opts Options) []string {
	var result []string
	var prevLine string
	lineCount := 0
	isDuplicatePrinted := false

	for i, currentLine := range lines {
		// Сравниваем текущую строку с предыдущей
		isSameLine := CompareStrings(currentLine, prevLine, opts.IgnoreCase, opts.NumFields, opts.NumChars)

		switch opts.Mode {
		case "default":
			// Если "default" - выводим строку, если она не совпадает с предыдущей
			if !isSameLine {
				result = append(result, currentLine)
			}

		case "count":
			// В режиме "count" подсчитываем количество повторений каждой строки
			if isSameLine {
				lineCount++
			} else {
				if lineCount > 0 {
					result = append(result, fmt.Sprintf("%d %s", lineCount, prevLine))
				}
				lineCount = 1
				prevLine = currentLine
			}

		case "duplicate":
			// При "duplicate" выводим только повторяющиеся строки
			if isSameLine {
				lineCount++
				if lineCount > 1 && !isDuplicatePrinted {
					result = append(result, currentLine)
					isDuplicatePrinted = true
				}
			} else {
				lineCount = 1
				prevLine = currentLine
				isDuplicatePrinted = false
			}

		case "unique":
			// В режиме "unique" выводим только уникальные строки
			if !isSameLine && lineCount == 0 && i != 0 {
				result = append(result, prevLine)
			}
			if !isSameLine {
				lineCount = 0
				prevLine = currentLine
			} else {
				lineCount++
			}
		}

		prevLine = currentLine
	}

	// Обрабатываем последнюю строку
	if opts.Mode == "count" && lineCount > 0 {
		result = append(result, fmt.Sprintf("%d %s", lineCount, prevLine))
	} else if opts.Mode == "unique" && lineCount == 0 {
		result = append(result, prevLine)
	}

	return result
}
