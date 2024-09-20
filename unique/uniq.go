package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"awesomeProject5/handleflags"
	"awesomeProject5/uniq"
)

func main() {

	var flags handleflags.Flags
	var err error
	// Вызываем функцию handleFlags для обработки флагов командной строки
	flags, err = handleflags.HandleFlags()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Чтение строк
	var reader io.Reader
	file, err := os.Open(flags.InputFile)
	if err == nil {
		reader = file
		defer file.Close()
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
		Mode:       flags.Mode,
		IgnoreCase: flags.IgnoreCase,
		NumFields:  flags.NumFields,
		NumChars:   flags.NumChars,
	}

	// Обрабатываем строки с помощью новой функции
	result := uniq.ProcessStrings(lines, opts)

	// Записываем результат
	var writer io.Writer
	if flags.OutputFile != "" {
		file, err := os.Create(flags.OutputFile)
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
