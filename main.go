package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/wipdev-tech/go-pokedex/internal/pokeapi"
	"github.com/wipdev-tech/go-pokedex/internal/pokecache"
)

type cliCmd struct {
	name        string
	description string
	callback    func(*cmdConfig, string) error
}

type cmdConfig struct {
	next     string
	previous string
	cache    pokecache.Cache
	pokedex  map[string]pokeapi.Pokemon
}

var cfg = cmdConfig{
	next:    "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
	cache:   pokecache.NewCache(5 * time.Minute),
	pokedex: map[string]pokeapi.Pokemon{},
}

var cmds = map[string]cliCmd{
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    cmdHelp,
	},
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    cmdExit,
	},
	"map": {
		name:        "map",
		description: "Get next 20 locations",
		callback:    cmdMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Get previous 20 locations",
		callback:    cmdMapB,
	},
	"explore": {
		name:        "explore",
		description: "Explore an area",
		callback:    cmdExplore,
	},
	"catch": {
		name:        "catch",
		description: "Try to catch a Pokemon",
		callback:    cmdCatch,
	},
}

func main() {
	fmt.Print("◓ > ")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		userInput := strings.Fields(s.Text())
		userCmd := userInput[0]
		var arg string
		if len(userInput) > 1 {
			arg = userInput[1]
		}
		if cliCmd, ok := cmds[userCmd]; ok {
			err := cliCmd.callback(&cfg, arg)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Invalid command: use `help` if you're stuck.")
		}
		fmt.Print("◓ > ")
	}
}

func cmdHelp(cfg *cmdConfig, area string) error {
	fmt.Println(`
Welcome to the Pokedex!

Usage:
help  Displays a help message
map   Shows the next 20 map locations
mapb  Shows the previous 20 map locations
exit  Exit the Pokedex`)
	fmt.Println()
	return nil
}

func cmdExit(cfg *cmdConfig, area string) error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}

func cmdMap(cfg *cmdConfig, area string) error {
	url := cfg.next
	if cfg.next == "" {
		return fmt.Errorf("already at the end")
	}

	var locations pokeapi.LocationData
	if cached, ok := cfg.cache.Get(url); ok {
		err := json.Unmarshal(cached, &locations)
		if err != nil {
			return err
		}
	} else {
		result, err := pokeapi.GetLocations(url)
		if err != nil {
			return err
		}
		locations = result

		resultBytes, err := json.Marshal(locations)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, resultBytes)
	}

	cfg.next = locations.Next
	cfg.previous = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func cmdMapB(cfg *cmdConfig, area string) error {
	url := cfg.previous
	if url == "" {
		return fmt.Errorf("already at the start")
	}

	var locations pokeapi.LocationData
	if cached, ok := cfg.cache.Get(url); ok {
		err := json.Unmarshal(cached, &locations)
		if err != nil {
			return err
		}
	} else {
		result, err := pokeapi.GetLocations(url)
		if err != nil {
			return err
		}
		locations = result

		resultBytes, err := json.Marshal(locations)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, resultBytes)
	}

	cfg.next = locations.Next
	cfg.previous = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func cmdExplore(cfg *cmdConfig, area string) error {
	if area == "" {
		return fmt.Errorf("No area given")
	}
	url := "https://pokeapi.co/api/v2/location-area/" + area
	fmt.Printf("Exploring %s...\n", area)

	var exploreData pokeapi.ExploreData
	if cached, ok := cfg.cache.Get(url); ok {
		err := json.Unmarshal(cached, &exploreData)
		if err != nil {
			return err
		}
	} else {
		result, err := pokeapi.GetExplore(url)
		if err != nil {
			return err
		}
		exploreData = result

		resultBytes, err := json.Marshal(exploreData)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, resultBytes)
	}

	fmt.Println("Found Pokemon:")
	for _, pokemonEncounter := range exploreData.PokemonEncounters {
		fmt.Println("-", pokemonEncounter.Pokemon.Name)
	}

	return nil
}

func cmdCatch(cfg *cmdConfig, name string) error {
	if name == "" {
		return fmt.Errorf("No Pokemon given")
	}
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	var pokemon pokeapi.Pokemon
	if cached, ok := cfg.cache.Get(url); ok {
		err := json.Unmarshal(cached, &pokemon)
		if err != nil {
			return err
		}
	} else {
		result, err := pokeapi.GetPokemon(url)
		if err != nil {
			return err
		}
		pokemon = result

		resultBytes, err := json.Marshal(pokemon)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, resultBytes)
	}

	if catch(pokemon.BaseExperience) {
		fmt.Printf("%s was caught!\n", name)
		cfg.pokedex[name] = pokemon
		fmt.Println(cfg.pokedex)
		return nil
	}

	fmt.Printf("%s escaped!\n", name)
	return nil
}

func catch(baseExp int) bool {
	const minBaseExp = 64
	const maxBaseExp = 635
	randVal := rand.Intn(maxBaseExp-minBaseExp) + minBaseExp
	return randVal > baseExp
}
