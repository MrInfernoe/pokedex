package repl_test

import (
	"testing"
	"pokedex/internal/repl"
	// "bufio"
	"os"
	"os/exec"
	// "fmt"
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
		got := repl.CleanInput(tc.input)
		for i := range got {
			word := got[i]
			expected := tc.output[i]
			if word != expected {
				t.Errorf("output: %q does not match expected: %q", got, tc.output)
			}
		}
	}
}

// Andrew Gerrand https://go.dev/talks/2014/testing.slide#23

func TestCommandExit(t *testing.T) {
    if os.Getenv("BE_EXITER") == "1" {
        repl.CommandExit()
        return
    }
    cmd := exec.Command(os.Args[0], "-test.run=TestCommandExit")
    cmd.Env = append(os.Environ(), "BE_EXITER=1")
    err := cmd.Run()
    if e, ok := err.(*exec.ExitError); ok && !e.Success() {
        return
    }
    t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestCommandHelp(t *testing.T) {
	repl.CommandHelp()
}