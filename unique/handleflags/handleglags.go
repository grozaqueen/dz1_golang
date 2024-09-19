package handleflags

import (
	"flag"
	"fmt"
)

// Функция handleFlags обрабатывает флаги командной строки
func HandleFlags() (string, string, string, bool, int, int, error) {
	countFlag := flag.Bool("c", false, "Подсчитать количество встречаний строки") // -с для подсчета количества повторений строки
	duplicateFlag := flag.Bool("d", false, "Вывести только повторяющиеся строки") // -d для вывода только повторяющихся строк
	uniqueFlag := flag.Bool("u", false, "Вывести только уникальные строки")       // -u для вывода только уникальных строк
	ignoreCaseFlag := flag.Bool("i", false, "Игнорировать регистр букв")          // -i для игнорирования регистра букв при сравнении строк
	numFieldsFlag := flag.Int("f", 0, "Не учитывать первые num_fields полей")     // -f для указания количества полей, которые нужно игнорировать в начале каждой строки
	numCharsFlag := flag.Int("s", 0, "Не учитывать первые num_chars символов")    // -s для указания количества символов, которые нужно игнорировать в начале каждой строки

	flag.Parse() // Парсим флаги командной строки

	// Проверяем, не были ли переданы взаимоисключающие флаги -c, -d, -u
	if (*countFlag && *duplicateFlag) || (*countFlag && *uniqueFlag) || (*duplicateFlag && *uniqueFlag) {
		return "", "", "", false, 0, 0, fmt.Errorf("ошибка: Флаги -c, -d и -u взаимоисключающие")
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
		return "", "", "", false, 0, 0, fmt.Errorf("ошибка: Неизвестные флаги: -%s", invalidFlags)
	}

	// Определяем имена входного (inputFile) и выходного (outputFile) файлов на основе переданных аргументов
	var inputFile, outputFile string
	switch flag.NArg() { // NArg — количество аргументов, оставшихся после обработки флагов
	case 0:
		inputFile = "" // Если не передан ни один аргумент, то - чтение из стандартного ввода (stdin)
	case 1:
		inputFile = flag.Arg(0) //  имя входного файла (inputFile)
	case 2:
		inputFile = flag.Arg(0)  //  первый - имя входного файла (inputFile)
		outputFile = flag.Arg(1) // второй - имя выходного файла (outputFile)
	default:
		return "", "", "", false, 0, 0, fmt.Errorf("ошибка: Слишком много аргументов") // Возвращаем ошибку
	}

	// определяем mode на основе переданных флагов
	var mode string
	if *countFlag {
		mode = "count"
	} else if *duplicateFlag {
		mode = "duplicate"
	} else if *uniqueFlag {
		mode = "unique"
	} else {
		mode = "default"
	}

	return mode, inputFile, outputFile, *ignoreCaseFlag, *numFieldsFlag, *numCharsFlag, nil
}
