package uniq

import (
	"fmt"
	"strings"
)

const (
	SpaceRune     = ' '
	TabRune       = '\t'
	DefaultMode   = "default"
	CountMode     = "count"
	DuplicateMode = "duplicate"
	UniqueMode    = "unique"
	SymbolEmpty   = ""
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
	field := SymbolEmpty
	inField := false

	for _, r := range s {
		if r == SpaceRune || r == TabRune { // Пробел или табуляция считаются разделителями
			if inField { // Если мы были внутри поля, то конец поля
				fields = append(fields, field)
				field = SymbolEmpty
			}
			inField = false // Выходим из поля
		}
		inField = true
		field += string(r) // Добавляем символ к полю
	}

	// Добавляем последнее поле, если строка не закончилась разделителем
	if inField {
		fields = append(fields, field)
	}

	return fields
}

// функция сравнения пары строк с учетом введенных параметров
func CompareStrings(str1, str2 string, opts Options) bool {
	// Игнорируем регистр, если установлен флаг
	if opts.IgnoreCase {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	// Если нужно игнорировать numFields полей
	if opts.NumFields > 0 {
		fields1 := splitIntoFields(str1)
		fields2 := splitIntoFields(str2)

		// Обрезаем поля для str1
		if len(fields1) > opts.NumFields {
			str1 = SymbolEmpty
			for i := opts.NumFields; i < len(fields1); i++ {
				if i > opts.NumFields {
					str1 += string(SpaceRune) // Вставляем пробелы между полями
				}
				str1 += fields1[i]
			}
		}

		// Аналогично для str2
		if len(fields2) > opts.NumFields {
			str2 = SymbolEmpty
			for i := opts.NumFields; i < len(fields2); i++ {
				if i > opts.NumFields {
					str2 += string(SpaceRune)
				}
				str2 += fields2[i]
			}
		}
	}

	// Игнорируем numChars символов путем среза
	if opts.NumChars > 0 {
		if len(str1) >= opts.NumChars {
			str1 = str1[opts.NumChars:]
		}
		if len(str2) >= opts.NumChars {
			str2 = str2[opts.NumChars:]
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

	processPrevLine := func() {
		if lineCount > 0 && opts.Mode == CountMode {
			result = append(result, fmt.Sprintf("%d %s", lineCount, prevLine))
		} else if lineCount == 0 && opts.Mode == UniqueMode {
			result = append(result, prevLine)
		}
	}

	for i, currentLine := range lines {
		isSameLine := CompareStrings(currentLine, prevLine, opts)

		switch opts.Mode {
		case DefaultMode:
			if !isSameLine {
				result = append(result, currentLine)
			}
		case CountMode:
			if isSameLine {
				lineCount++
			} else {
				processPrevLine()
				lineCount = 1
			}
		case DuplicateMode:
			if isSameLine {
				lineCount++
				if lineCount > 1 && !isDuplicatePrinted {
					result = append(result, currentLine)
					isDuplicatePrinted = true
				}
			} else {
				lineCount = 1
				isDuplicatePrinted = false
			}
		case UniqueMode:
			if !isSameLine && lineCount == 0 && i != 0 {
				result = append(result, prevLine)
			}
			if !isSameLine {
				lineCount = 0
			} else {
				lineCount++
			}
		}

		prevLine = currentLine
	}

	processPrevLine()

	return result
}
