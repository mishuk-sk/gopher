package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadQuiz(t *testing.T) {

	tests := []struct {
		input      string
		expected   []problem
		shouldFail bool
	}{
		{
			input:      "1+2,3\n",
			expected:   []problem{{"1+2", "3"}},
			shouldFail: false,
		},
		{
			input:      "1+3,    4",
			expected:   []problem{{"1+3", "4"}},
			shouldFail: false,
		},
		{
			input:      "1+1,2, 20\n2+3, 5, 7",
			expected:   []problem{},
			shouldFail: true,
		},
		{
			input:      "1+1,2\n2+3, 5",
			expected:   []problem{{"1+1", "2"}, {"2+3", "5"}},
			shouldFail: false,
		},
	}
	for _, test := range tests {
		r := strings.NewReader(test.input)
		problems := readQuiz(r)
		if !reflect.DeepEqual(problems, test.expected) {
			t.Errorf("For input '%s' expected '%s' but got '%s'", test.input, test.expected, problems)
		}
	}
}
