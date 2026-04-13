package repl

import (
	"strings"
	"fmt"
	"os"
	"net/http"
	"io"
	"encoding/json"
)

type Config struct {
	Next		string
	Previous	string
}


//split the user's input into "words" based on whitespace. It should also lowercase the input and trim any leading or trailing whitespace.
func CleanInput(text string) []string {
	lowered := strings.ToLower(text)
	splitted := strings.Fields(lowered)
	return splitted
}
// registry command to Exit
func CommandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
// registry command displays help
func CommandHelp(c *Config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Display a help message\nexit: Exit the Pokedex")
	return nil
}


type LocationArea struct {
	Count		int					`json:"count"`
	Next		string				`json:"next"`
	Previous	string				`json:"previous"`
	Results		[]map[string]string	`json:"results"`
}

// displays the names of 20 location areas in the Pokemon world. Each subsequent call to map should display the next 20 locations
func CommandMap(c *Config) error {
	var link string
	if c.Next == "" {
		link = "https://pokeapi.co/api/v2/location-area/"
	} else {
		link = c.Next
	}
	res, err := http.Get(link)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// fmt.Printf("body: %s\n", body)

	var LocA LocationArea
	err = json.Unmarshal(body, &LocA)
	if err != nil {
		return err
	}
	// fmt.Println(LocA)
	for _, location := range LocA.Results {
		fmt.Printf("%s\n", location["name"])
	}
	c.Next = LocA.Next
	c.Previous = LocA.Previous

	return nil
}

func CommandMapBack(c *Config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	c.Next = c.Previous
	err := CommandMap(c)
	return err
}

type CliCommand struct {
	name		string
	description	string
	Callback	func(c *Config) error
}
// ---change to function to avoid init cycle
func GetRegistry() map[string]CliCommand {
	registry := map[string]CliCommand{
		"help": {
			name:			"help",
			description:	"Display a help message",
			Callback:		CommandHelp,
		},
		"exit": {
			name:			"exit",
			description:	"Exit the pokedex",
			Callback:		CommandExit,
		},
		"map": {
			name:			"map",
			description:	"Display the next 20 location names",
			Callback:		CommandMap,
		},
		"mapb": {
			name:			"mapb",
			description: 	"Display the previous 20 location names",
			Callback:		CommandMapBack,
		},
	}
	return registry
}

