package langShell_test

import (
	"testing"
	"pokedex/internal/langShell"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
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
		got := langShell.CleanInput(tc.input)
		for i := range got {
			word := got[i]
			expected := tc.output[i]
			if word != expected {
				t.Errorf("output: %q does not match expected: %q", got, tc.output)
			}
		}
	}
}