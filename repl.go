package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/Enri258/pokedexcli/internal/pokecache"
)

type config struct {
	commands            map[string]cliCommand
	nextLocationURL     string
	previousLocationURL string
	cache               *pokecache.Cache
	pokedex             map[string]Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range cfg.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(cfg *config, args ...string) error {
	if cfg.nextLocationURL == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	data, err := getLocationAreas(
		cfg.nextLocationURL,
		cfg.cache,
	)
	if err != nil {
		return err
	}

	if data.Next != nil {
		cfg.nextLocationURL = *data.Next
	} else {
		cfg.nextLocationURL = ""
	}

	if data.Previous != nil {
		cfg.previousLocationURL = *data.Previous
	} else {
		cfg.previousLocationURL = ""
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapBack(cfg *config, args ...string) error {
	if cfg.previousLocationURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	data, err := getLocationAreas(
		cfg.previousLocationURL,
		cfg.cache,
	)
	if err != nil {
		return err
	}

	if data.Next != nil {
		cfg.nextLocationURL = *data.Next
	} else {
		cfg.nextLocationURL = ""
	}

	if data.Previous != nil {
		cfg.previousLocationURL = *data.Previous
	} else {
		cfg.previousLocationURL = ""
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <area_name>")
	}

	areaName := args[0]

	fmt.Printf("Exploring %s...\n", areaName)

	area, err := getLocationArea(areaName, cfg.cache)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")

	for _, encounter := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon_name>")
	}

	pokemonName := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := getPokemon(pokemonName, cfg.cache)
	if err != nil {
		return err
	}

	const catchThreshold = 50

	roll := rand.Intn(pokemon.BaseExperience + 1)

	if roll > catchThreshold {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemon.Name)
	fmt.Println("You may now inspect it with the inspect command.")

	cfg.pokedex[pokemon.Name] = pokemon

	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon_name>")
	}

	pokemonName := args[0]

	pokemon, caught := cfg.pokedex[pokemonName]
	if !caught {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("  - %s\n", pokemonType.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	fmt.Println("Your Pokedex:")

	for name := range cfg.pokedex {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			return
		}

		words := cleanInput(scanner.Text())

		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := words[1:]

		command, exists := cfg.commands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(cfg, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}
