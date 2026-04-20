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
	"math/rand"
	// "math"
)



type Config struct {
	Next			string
	Previous		string
	LinkCache		*pokecache.Cache
	HelpOrder		[]string
	Pokedex			map[string]Pokemon
}
func NewConfig() *Config {
	var config Config
	config.HelpOrder = []string{"help", "map", "mapb", "explore", "catch", "pokedex", "inspect", "exit"}
	config.LinkCache = pokecache.NewCache(5*time.Second)
	config.Pokedex = map[string]Pokemon{}
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

	if string(body) == "Not Found" {
		fmt.Printf("The location %s could not be found.\n", locationName)
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



//	Add a catch command. It takes the name of a Pokemon as an argument. Example:

// Pokedex > catch pikachu
// Throwing a Pokeball at pikachu...
// pikachu escaped!
// Pokedex > catch pikachu
// Throwing a Pokeball at pikachu...
// pikachu was caught!

//	Be sure to print the Throwing a Pokeball at <pokemon>... message before 
// 	determining if the Pokemon was caught or not.

//	Use the Pokemon endpoint to get information about a Pokemon by name.

//	Give the user a chance to catch the Pokemon using the math/rand package.

//	You can use the pokemon's "base experience" to determine the chance of 
// 	catching it. The higher the base experience, the harder it should be to catch.

//	Once the Pokemon is caught, add it to the user's Pokedex. I used a 
// 	map[string]Pokemon to keep track of caught Pokemon.

//	Test the catch command manually - make sure you can actually catch a Pokemon 
// 	within a reasonable number of tries.

// name, height, weight, stats and type(s)
type Pokemon struct {
	name	string
	height	int
	weight	int
	stats	map[string]int
	types	[]string
}

type PokemonData struct {
	BaseExperience	int		`json:"base_experience"`
	Height			int		`json:"height"`
	Stats			[]struct {
		BaseStat		int		`json:"base_stat"`
		Stat			struct {
			Name			string	`json:"name"`
		}						`json:"stat"`
	}						`json:"stats"`
	Types			[]struct {
		Type			struct {
			Name			string	`json:"name"`
		}						`json:"type"`
	}						`json:"types"`
	Weight			int		`json:"weight"`
}

func PDtoPokemon(name string, PD PokemonData) Pokemon {
	newP := Pokemon{}
	newP.name = name
	newP.height = PD.Height
	newP.weight = PD.Weight

	newP.stats = map[string]int{}
	for _, stat := range PD.Stats {
		newP.stats[stat.Stat.Name] = stat.BaseStat
	}

	newP.types = []string{}
	for _, thisType := range PD.Types {
		newP.types = append(newP.types, thisType.Type.Name)
	}
	
	return newP
}

func CommandCatch(c *Config, pokemonName string) error {
	// chance = (highest-this)/highest

	// make link
	// look in cache for body
	// if no entry
	// 	get response
	//	read response body
	// 	defer Close
	// 	add to cache
	//	Unmarshal into PokemonData
	// 	
	
	baseExperienceLink := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	var body []byte
	if baseExperience, exists := c.LinkCache.Get(baseExperienceLink); exists {
		body = baseExperience
		} else {
			res, err := http.Get(baseExperienceLink)
			if err != nil {
				return err
			}
			body, err = io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			defer res.Body.Close()
			c.LinkCache.Add(baseExperienceLink, body)
		}
		
	if string(body) == "Not Found" {
		fmt.Printf("The pokemon named %s could not be found.\n", pokemonName)
		return nil
	}
		
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	var PD PokemonData
	err := json.Unmarshal(body, &PD)
	if err != nil {
		return err
	}
	// chance for highestExp pokemon (635xp) = 10%
	// chance for lowestExp pokemon (36xp) = 90%
	// linear interpolation: chance = 100 - 80/599(baseExperience + 38.875)

	catchChance := 100.0 - 80.0/599.0*(float64(PD.BaseExperience) + 38.0)
	random100 := rand.Intn(100)
	if random100 <= int(catchChance) {
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
		c.Pokedex[pokemonName] = PDtoPokemon(pokemonName, PD)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}



func CommandInspect(c *Config, pokemonName string) error {
	p, exists := c.Pokedex[pokemonName]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", p.name)
	fmt.Printf("Height: %d\n", p.height)
	fmt.Printf("Weight: %d\n", p.weight)
	fmt.Println("Stats:")

	fmt.Printf("\t-hp: %d\n", p.stats["hp"])
	fmt.Printf("\t-attack: %d\n", p.stats["attack"])
	fmt.Printf("\t-defense: %d\n", p.stats["defense"])
	fmt.Printf("\t-special-attack: %d\n", p.stats["special-attack"])
	fmt.Printf("\t-special-defense: %d\n", p.stats["special-defense"])
	fmt.Printf("\t-speed: %d\n", p.stats["speed"])
	
	fmt.Println("Types:")
	for _, thisType := range p.types {
		fmt.Printf("\t- %s\n", thisType)
	}

	return nil
}



func CommandPokedex(c *Config, empty string) error {
	fmt.Println("Your Pokedex:")
	for pokemonName, _ := range c.Pokedex {
		fmt.Printf("\t- %s\n", pokemonName)
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
		"catch": {
			name:			"catch",
			description:	"Attempt to catch a pokemon",
			Callback:		CommandCatch,
		},
		"inspect": {
			name:			"inspect",
			description:	"Displays data for a caught pokemon",
			Callback:		CommandInspect,
		},
		"pokedex": {
			name:			"pokedex",
			description:	"Display caught pokemon",
			Callback:		CommandPokedex,
		},
	}
	return registry
}