package main

import (
	"time"

	"github.com/Enri258/pokedexcli/internal/pokecache"
)

func main() {
	cfg := &config{
		nextLocationURL: "https://pokeapi.co/api/v2/location-area/",
		cache:           pokecache.NewCache(5 * time.Second),
		pokedex:         make(map[string]Pokemon),
		commands: map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
			"map": {
				name:        "map",
				description: "Display the next 20 location areas",
				callback:    commandMap,
			},
			"mapb": {
				name:        "mapb",
				description: "Display the previous 20 location areas",
				callback:    commandMapBack,
			},
			"explore": {
				name:        "explore",
				description: "Explore a location area",
				callback:    commandExplore,
			},
			"catch": {
				name:        "catch",
				description: "Catch a Pokemon",
				callback:    commandCatch,
			},
			"inspect": {
				name:        "inspect",
				description: "Inspect a caught Pokemon",
				callback:    commandInspect,
			},
			"pokedex": {
				name:        "pokedex",
				description: "List all caught Pokemon",
				callback:    commandPokedex,
			},
		},
	}

	startRepl(cfg)
}
