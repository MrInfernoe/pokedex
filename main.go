package main

import (
	"bufio"
	"os"
	"fmt"
	repl "pokedex/internal/repl"
)

/*
A REPL is only useful if it does something! Our REPL will work using the concept of "commands". A command is a single word that maps to an action.

We're going to support two commands in this step:

    help: prints a help message describing how to use the REPL
    exit: exits the program

Assignment

-    Remove your logic that prints the first word (the command) back to the user
    Add a callback for the exit command. Commands in our REPL are just callback functions with no arguments, but that return an error. For example:

func commandExit() error

This function should print Closing the Pokedex... Goodbye! then immediately exit the program. I used os.Exit(0).

    Create a "registry" of commands. This will give us a nice abstraction for managing the many commands we'll be adding. I created a struct type that describes a command:

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

Then I created a map of supported commands:

map[string]cliCommand{
    "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    },
}

    Register the exit command. Update your REPL loop to use the "command" the user typed in to look up the callback function in the registry. If the command is found, call the callback (and print any errors that are returned). If there isn't a handler, just print Unknown command.
    Test your program (obviously).
    Add a help command, its callback, and register it. It should print:

Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

You can dynamically generate the "usage" section by iterating over my registry of commands. That way the help command will always be up-to-date with the available commands.

    Test your code again manually.

*/

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

func main() {
registry := map[string]cliCommand{
	"exit": {
		name:			"exit",
		description:	"Exit the pokedex",
		callback:		repl.CommandExit,
	},
	"help": {
		name:			"help",
		description:	"Display a help message",
		callback:		repl.CommandHelp,
	},
}

	scanner := bufio.NewScanner(os.Stdin)
	for true {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			fmt.Println(scanner.Err())
			break
		}
		userText := scanner.Text()
		cleanText := repl.CleanInput(userText)
		firstWord := cleanText[0]
		callFunc := registry[firstWord].callback
		if callFunc != nil {
			callFunc()
		} else {
			fmt.Println(fmt.Errorf("Unknown command"))
		}
	}
}