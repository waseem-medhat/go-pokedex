package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCmd struct {
	name        string
	description string
	callback    func() error
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
}

func main() {
	fmt.Print("Pokedex > ")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		userCmd := s.Text()
		cmds[userCmd].callback()
		fmt.Print("Pokedex > ")
	}
}

func cmdHelp() error {
	fmt.Println(`
Welcome to the Pokedex!

Usage:
help  Displays a help message
exit  Exit the Pokedex`)
    fmt.Println()
	return nil
}

func cmdExit() error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}
