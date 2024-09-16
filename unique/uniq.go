package main

import (
	"awesomeProject5/uniq"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

// Функция handleFlags обрабатывает флаги командной строки
func handleFlags() (string, string, string, bool, int, int, error) {
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

func main() {
	// Вызываем функцию handleFlags для обработки флагов командной строки
	mode, inputFile, outputFile, ignoreCase, numFields, numChars, err := handleFlags()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Чтение строк
	var reader io.Reader
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Println("Ошибка при открытии файла:", err)
			return
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	// Чтение строк из reader
	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Настройка опций
	opts := uniq.Options{
		Mode:       mode,
		IgnoreCase: ignoreCase,
		NumFields:  numFields,
		NumChars:   numChars,
	}

	// Обрабатываем строки с помощью новой функции
	result := uniq.ProcessStrings(lines, opts)

	// Записываем результат
	var writer io.Writer
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Ошибка при создании файла:", err)
			return
		}
		defer file.Close()
		writer = file
	} else {
		writer = os.Stdout
	}

	// Выводим результат
	for _, line := range result {
		fmt.Fprintln(writer, line)
	}
}
