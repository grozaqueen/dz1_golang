package main

import (
	"awesomeProject5/uniq"
	"flag"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_processStrings_Success(t *testing.T) {
	tests := []struct {
		name       string
		input      []string
		mode       string
		ignoreCase bool
		numFields  int
		numChars   int
		want       []string
	}{
		{
			name:       "default mode, no options",
			input:      []string{"a", "b", "a"},
			mode:       "default",
			ignoreCase: false,
			numFields:  0,
			numChars:   0,
			want:       []string{"a", "b", "a"},
		},
		{
			name:       "count mode",
			input:      []string{"a", "b", "a"},
			mode:       "count",
			ignoreCase: false,
			numFields:  0,
			numChars:   0,
			want:       []string{"1 a", "1 b", "1 a"},
		},
		{
			name:       "duplicate mode",
			input:      []string{"a", "b", "a", "a"},
			mode:       "duplicate",
			ignoreCase: false,
			numFields:  0,
			numChars:   0,
			want:       []string{"a"},
		},
		{
			name:       "unique mode",
			input:      []string{"a", "b", "a"},
			mode:       "unique",
			ignoreCase: false,
			numFields:  0,
			numChars:   0,
			want:       []string{"a", "b", "a"},
		},
		{
			name:       "ignore case",
			input:      []string{"a", "A", "b"},
			mode:       "default",
			ignoreCase: true,
			numFields:  0,
			numChars:   0,
			want:       []string{"a", "b"},
		},
		{
			name:       "ignore fields",
			input:      []string{"a b c", "a d e"},
			mode:       "default",
			ignoreCase: false,
			numFields:  1,
			numChars:   0,
			want:       []string{"a b c", "a d e"},
		},
		{
			name:       "ignore chars",
			input:      []string{"abc", "abd"},
			mode:       "default",
			ignoreCase: false,
			numFields:  0,
			numChars:   1,
			want:       []string{"abc", "abd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := uniq.Options{
				Mode:       tt.mode,
				IgnoreCase: tt.ignoreCase,
				NumFields:  tt.numFields,
				NumChars:   tt.numChars,
			}
			got := uniq.ProcessStrings(tt.input, opts)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCompareStrings_Success(t *testing.T) {
	tests := []struct {
		name       string
		str1       string
		str2       string
		ignoreCase bool
		numFields  int
		numChars   int
		want       bool
	}{
		{
			name:       "EqualStrings",
			str1:       "hello",
			str2:       "hello",
			ignoreCase: false,
			numFields:  0,
			numChars:   0,
			want:       true,
		},
		{
			name:       "DifferentStrings",
			str1:       "hello",
			str2:       "world",
			ignoreCase: false,
			numFields:  0,
			numChars:   0,
			want:       false,
		},
		{
			name:       "IgnoreCase_EqualStrings",
			str1:       "hello",
			str2:       "HELLO",
			ignoreCase: true,
			numFields:  0,
			numChars:   0,
			want:       true,
		},
		{
			name:       "IgnoreCase_DifferentStrings",
			str1:       "hello",
			str2:       "WORLD",
			ignoreCase: true,
			numFields:  0,
			numChars:   0,
			want:       false,
		},
		{
			name:       "NumFields_EqualStrings",
			str1:       "1 hello world",
			str2:       "2 hello world",
			ignoreCase: false,
			numFields:  1,
			numChars:   0,
			want:       true,
		},
		{
			name:       "NumFields_DifferentStrings",
			str1:       "1 hello world",
			str2:       "2 world hello",
			ignoreCase: false,
			numFields:  1,
			numChars:   0,
			want:       false,
		},
		{
			name:       "NumChars_EqualStrings",
			str1:       "abchello",
			str2:       "abclowo",
			ignoreCase: false,
			numFields:  0,
			numChars:   3,
			want:       false,
		},
		{
			name:       "NumChars_DifferentStrings",
			str1:       "abchello",
			str2:       "abcworld",
			ignoreCase: false,
			numFields:  0,
			numChars:   3,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uniq.CompareStrings(tt.str1, tt.str2, tt.ignoreCase, tt.numFields, tt.numChars)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestHandleFlags_Errors(t *testing.T) {
	tests := []struct {
		name               string
		flags              []string
		expectedMode       string
		expectedInputFile  string
		expectedOutputFile string
		expectedIgnoreCase bool
		expectedNumFields  int
		expectedNumChars   int
		expectedError      error
	}{
		{
			name:          "conflicting_flags_cd",
			flags:         []string{"-c", "-d"},
			expectedError: fmt.Errorf("ошибка: Флаги -c, -d и -u взаимоисключающие"),
		},
		{
			name:          "conflicting_flags_cu",
			flags:         []string{"-c", "-u"},
			expectedError: fmt.Errorf("ошибка: Флаги -c, -d и -u взаимоисключающие"),
		},
		{
			name:          "conflicting_flags_du",
			flags:         []string{"-d", "-u"},
			expectedError: fmt.Errorf("ошибка: Флаги -c, -d и -u взаимоисключающие"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем старые аргументы командной строки
			oldArgs := os.Args
			defer func() {
				// Восстанавливаем старые аргументы командной строки
				os.Args = oldArgs
			}()

			// Сбрасываем состояние пакета flag перед каждым тестом
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Устанавливаем новые аргументы командной строки для теста
			os.Args = append([]string{"cmd"}, tt.flags...)

			// Вызываем тестируемую функцию
			mode, inputFile, outputFile, ignoreCase, numFields, numChars, err := handleFlags()

			// Проверяем наличие ошибки
			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			// Проверяем остальные значения, только если не было ошибок
			require.NoError(t, err)
			require.Equal(t, tt.expectedMode, mode)
			require.Equal(t, tt.expectedInputFile, inputFile)
			require.Equal(t, tt.expectedOutputFile, outputFile)
			require.Equal(t, tt.expectedIgnoreCase, ignoreCase)
			require.Equal(t, tt.expectedNumFields, numFields)
			require.Equal(t, tt.expectedNumChars, numChars)
		})
	}
}
