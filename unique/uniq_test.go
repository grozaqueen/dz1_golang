package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"awesomeProject5/handleflags"
	"awesomeProject5/uniq"
	"github.com/stretchr/testify/require"
)

func Test_processStrings_Success(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
		name string
		str1 string
		str2 string
		opts uniq.Options
		want bool
	}{
		{
			name: "EqualStrings",
			str1: "hello",
			str2: "hello",
			opts: uniq.Options{
				IgnoreCase: false,
				NumFields:  0,
				NumChars:   0,
			},
			want: true,
		},
		{
			name: "DifferentStrings",
			str1: "hello",
			str2: "world",
			opts: uniq.Options{
				IgnoreCase: false,
				NumFields:  0,
				NumChars:   0,
			},
			want: false,
		},
		{
			name: "IgnoreCase_EqualStrings",
			str1: "hello",
			str2: "HELLO",
			opts: uniq.Options{
				IgnoreCase: true,
				NumFields:  0,
				NumChars:   0,
			},
			want: true,
		},
		{
			name: "IgnoreCase_DifferentStrings",
			str1: "hello",
			str2: "WORLD",
			opts: uniq.Options{
				IgnoreCase: true,
				NumFields:  0,
				NumChars:   0,
			},
			want: false,
		},
		{
			name: "NumFields_EqualStrings",
			str1: "1 hello world",
			str2: "2 hello world",
			opts: uniq.Options{
				IgnoreCase: false,
				NumFields:  1,
				NumChars:   0,
			},
			want: true,
		},
		{
			name: "NumFields_DifferentStrings",
			str1: "1 hello world",
			str2: "2 world hello",
			opts: uniq.Options{
				IgnoreCase: false,
				NumFields:  1,
				NumChars:   0,
			},
			want: false,
		},
		{
			name: "NumChars_EqualStrings",
			str1: "abchello",
			str2: "abclowo",
			opts: uniq.Options{
				IgnoreCase: false,
				NumFields:  0,
				NumChars:   3,
			},
			want: false,
		},
		{
			name: "NumChars_DifferentStrings",
			str1: "abchello",
			str2: "abcworld",
			opts: uniq.Options{
				IgnoreCase: false,
				NumFields:  0,
				NumChars:   3,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uniq.CompareStrings(tt.str1, tt.str2, tt.opts)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestHandleFlags_Errors(t *testing.T) {
	tests := []struct {
		name          string
		flags         []string
		expectedError error
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
			name:          "too_many_args",
			flags:         []string{"-c", "input.txt", "output.txt", "extra.txt"},
			expectedError: fmt.Errorf("ошибка: Слишком много аргументов"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = append([]string{"cmd"}, tt.flags...)

			_, err := handleflags.HandleFlags()

			require.EqualError(t, err, tt.expectedError.Error())
		})
	}
}

func TestHandleFlags_Success(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		flags              []string
		expectedMode       string
		expectedIgnoreCase bool
		expectedNumFields  int
		expectedNumChars   int
	}{
		{
			name:         "valid_flags_count",
			flags:        []string{"-c"},
			expectedMode: "count",
		},
		{
			name:               "valid_flags_ignore_case",
			flags:              []string{"-c", "-i"},
			expectedMode:       "count",
			expectedIgnoreCase: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = append([]string{"cmd"}, tt.flags...)

			flags, err := handleflags.HandleFlags()
			require.NoError(t, err)

			require.Equal(t, tt.expectedMode, flags.Mode)
			require.Equal(t, tt.expectedIgnoreCase, flags.IgnoreCase)
			require.Equal(t, tt.expectedNumFields, flags.NumFields)
			require.Equal(t, tt.expectedNumChars, flags.NumChars)
		})
	}
}
