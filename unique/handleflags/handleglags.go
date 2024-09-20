package handleflags

import (
	"flag"
	"fmt"
)

// Структура для хранения всех флагов и параметров
type Flags struct {
	Mode       string // count, duplicate, unique, default
	InputFile  string // Входной файл
	OutputFile string // Выходной файл
	IgnoreCase bool   // Игнорировать регистр
	NumFields  int    // Количество полей для игнорирования
	NumChars   int    // Количество символов для игнорирования
}

// HandleFlags обрабатывает флаги командной строки и возвращает структуру Flags
func HandleFlags() (Flags, error) {
	// Определение флагов
	var (
		countFlag      bool
		duplicateFlag  bool
		uniqueFlag     bool
		ignoreCaseFlag bool
		numFieldsFlag  int
		numCharsFlag   int
	)

	flag.BoolVar(&countFlag, "c", false, "Подсчитать количество встречаний строки")
	flag.BoolVar(&duplicateFlag, "d", false, "Вывести только повторяющиеся строки")
	flag.BoolVar(&uniqueFlag, "u", false, "Вывести только уникальные строки")
	flag.BoolVar(&ignoreCaseFlag, "i", false, "Игнорировать регистр букв")
	flag.IntVar(&numFieldsFlag, "f", 0, "Не учитывать первые num_fields полей")
	flag.IntVar(&numCharsFlag, "s", 0, "Не учитывать первые num_chars символов")

	flag.Parse() // Парсим флаги командной строки

	// Проверяем, не были ли переданы взаимоисключающие флаги -c, -d, -u
	if (countFlag && duplicateFlag) || (countFlag && uniqueFlag) || (duplicateFlag && uniqueFlag) {
		return Flags{}, fmt.Errorf("ошибка: Флаги -c, -d и -u взаимоисключающие")
	}

	var invalidFlags []string
	validFlags := map[string]bool{
		"c": true,
		"d": true,
		"u": true,
		"i": true,
		"f": true,
		"s": true,
	}

	// Проверяем все переданные флаги
	flag.VisitAll(func(f *flag.Flag) {
		if !validFlags[f.Name] {
			invalidFlags = append(invalidFlags, f.Name)
		}
	})

	if len(invalidFlags) > 0 {
		return Flags{}, fmt.Errorf("ошибка: Неизвестные флаги: -%s", invalidFlags)
	}

	// Определяем имена входного (inputFile) и выходного (outputFile) файлов
	var inputFile, outputFile string
	switch flag.NArg() {
	case 0:
		inputFile = "" // Если не передан ни один аргумент — чтение из stdin
	case 1:
		inputFile = flag.Arg(0) // Имя входного файла
	case 2:
		inputFile = flag.Arg(0)  // Имя входного файла
		outputFile = flag.Arg(1) // Имя выходного файла
	default:
		return Flags{}, fmt.Errorf("ошибка: Слишком много аргументов")
	}

	// Определяем режим работы на основе флагов
	var mode string
	if countFlag {
		mode = "count"
	} else if duplicateFlag {
		mode = "duplicate"
	} else if uniqueFlag {
		mode = "unique"
	} else {
		mode = "default"
	}

	// Возвращаем структуру с результатами
	return Flags{
		Mode:       mode,
		InputFile:  inputFile,
		OutputFile: outputFile,
		IgnoreCase: ignoreCaseFlag,
		NumFields:  numFieldsFlag,
		NumChars:   numCharsFlag,
	}, nil
}
