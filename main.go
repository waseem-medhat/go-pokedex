package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/wipdev-tech/go-pokedex/internal/pokeapi"
	"github.com/wipdev-tech/go-pokedex/internal/pokecache"
)

type cliCmd struct {
	name        string
	description string
	callback    func(*cmdConfig) error
}

type cmdConfig struct {
	next     string
	previous string
	cache    pokecache.Cache
}

var cfg = cmdConfig{
	next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
	previous: "",
	cache:    pokecache.NewCache(15 * time.Second),
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
}

func main() {
	fmt.Print("◓ > ")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		userCmd := s.Text()
		if cliCmd, ok := cmds[userCmd]; ok {
			err := cliCmd.callback(&cfg)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Invalid command: use `help` if you're stuck.")
		}
		fmt.Print("◓ > ")
	}
}

func cmdHelp(cfg *cmdConfig) error {
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

func cmdExit(cfg *cmdConfig) error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}

func cmdMap(cfg *cmdConfig) error {
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

func cmdMapB(cfg *cmdConfig) error {
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
