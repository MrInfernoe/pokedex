package langShell

import (
	"strings"
)

//split the user's input into "words" based on whitespace. It should also lowercase the input and trim any leading or trailing whitespace.
func CleanInput(text string) []string {
	lowered := strings.ToLower(text)
	splitted := strings.Fields(lowered)
	return splitted
}