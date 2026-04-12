package repl

import (
	"strings"
	"fmt"
	"os"
)

//split the user's input into "words" based on whitespace. It should also lowercase the input and trim any leading or trailing whitespace.
func CleanInput(text string) []string {
	lowered := strings.ToLower(text)
	splitted := strings.Fields(lowered)
	return splitted
}
// registry command to Exit
func CommandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	// code 1 necessary for testing
	os.Exit(1)
	return nil
}
// registry command displays help
func CommandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}