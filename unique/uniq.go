package main

import (
	"bufio"
	"fmt"
	"io"
	
	"awesomeProject5/uniq"
	"awesomeProject5/handleflags"
)


func main() {
	// Вызываем функцию handleFlags для обработки флагов командной строки
	mode, inputFile, outputFile, ignoreCase, numFields, numChars, err := handleflags.HandleFlags()
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
