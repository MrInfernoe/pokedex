package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct{
		input string
		output []string
	}{
		{
			input: "helloworld",
			output: []string{"helloworld"},
		},{
			input: "hello world",
			output: []string{"hello", "world"},
		},{
			input: "hello/world",
			output: []string{"hello/world"},
		},{
			input: " hello   world    ",
			output: []string{"hello", "world"},
		},{
			input: "HellO WorlD",
			output: []string{"hello", "world"},
		},
	}

	for _, tc := range cases {
		got := cleanInput(tc.input)
		for i := range got {
			word := got[i]
			expected := tc.output[i]
			if word != expected {
				t.Errorf("output: %v does not match expected: %v", word, expected)
			}
		}
	}
}