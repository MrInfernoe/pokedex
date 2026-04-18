package repl

import (
	"strings"
	"fmt"
	"os"
	"net/http"
	"io"
	"encoding/json"
	"pokedex/internal/pokecache"
	"time"
	// "bufio"
)



type Config struct {
	Next			string
	Previous		string
	LinkCache		*pokecache.Cache
	HelpOrder		[]string
}
func NewConfig() *Config {
	var config Config
	config.HelpOrder = []string{"help", "map", "mapb", "explore", "exit"}
	config.LinkCache = pokecache.NewCache(5*time.Second)
	return &config
}



//split the user's input into "words" based on whitespace. It should also lowercase the input and trim any leading or trailing whitespace.
func CleanInput(text string) []string {
	lowered := strings.ToLower(text)
	splitted := strings.Fields(lowered)
	return splitted
}



// registry command to Exit
func CommandExit(c *Config, empty string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}



// registry command displays help
func CommandHelp(c *Config, empty string) error {
	registry := GetRegistry()
	order := c.HelpOrder
	if len(order) != len(registry) {
		return fmt.Errorf("missing command in ordering")
	}

	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, commandKey := range order {
		commandData := registry[commandKey]
		fmt.Printf("%-8s %s\n", commandData.name + ":", commandData.description)
	}
	return nil
}



// Map json Unmarshal
type LocationArea struct {
	// Count		int					`json:"count"`
	Next		string				`json:"next"`
	Previous	string				`json:"previous"`
	Results		[]map[string]string	`json:"results"`
}
// displays the names of 20 location areas in the Pokemon world. Each subsequent call to map should display the next 20 locations
func CommandMap(c *Config, empty string) error {
	var link string
	if c.Next == "" {
		link = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} else {
		link = c.Next
	}
	// fmt.Printf(link + "\n")

	var body []byte
	if getCache, ok := c.LinkCache.Get(link); ok {
		// fmt.Println("Restoring from cache")
		fmt.Printf("Restoring from cache:\n\n")
		body = getCache
	} else {
		// fmt.Println("Not restoring")
		res, err := http.Get(link)

		if err != nil {
			return err
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		// fmt.Printf("body: %s\n", body)

		c.LinkCache.Add(link, body)
	}

	var LocA LocationArea
	err := json.Unmarshal(body, &LocA)
	if err != nil {
		return err
	}
	// fmt.Println(LocA)
	for _, location := range LocA.Results {
		locationName := location["name"]
		locationName = strings.ReplaceAll(locationName, "-", " ")
		fmt.Printf("%s\n", locationName)
	}
	c.Next = LocA.Next
	c.Previous = LocA.Previous

	return nil
}
// back one page
func CommandMapBack(c *Config, empty string) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	c.Next = c.Previous
	err := CommandMap(c, "")
	return err
}



// Explore json Unmarshal
type AreaPokemon struct {
	Encounters	[]struct {
		Pokemon		struct{
			Name	string	`json:"name"`
			URL		string	`json:"url"`
		}	`json:"pokemon"`
	}	`json:"pokemon_encounters"`
}
func CommandExplore(c *Config, locationName string) error {
	if locationName == "" {
		return fmt.Errorf("no location given")
	}
	link := "https://pokeapi.co/api/v2/location-area/" + locationName

	var body []byte
	if pokemonData, exists := c.LinkCache.Get(link); exists {
		body = pokemonData
	} else {
		res, err := http.Get(link)
		if err != nil {
			return err
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		// print(string(body))
		c.LinkCache.Add(link, body)
	}
	var AD AreaPokemon
	err := json.Unmarshal(body, &AD)
	if err != nil {
		return err
	}

	for _, pokemonEncounter := range AD.Encounters {
		fmt.Printf("%s\n", pokemonEncounter.Pokemon.Name)
	}

	return nil
}



type CliCommand struct {
	name		string
	description	string
	Callback	func(c *Config, parameter string) error
}
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
		"explore": {
			name:			"explore",
			description:	"Display pokemon at location",
			Callback:		CommandExplore,
		},
	}
	return registry
}