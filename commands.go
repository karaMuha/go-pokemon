package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/karaMuha/pokedex/internal/config"
	"github.com/karaMuha/pokedex/internal/pokeapi"
	"github.com/karaMuha/pokedex/internal/pokedex"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config.Config) error
}

var commands = map[string]cliCommand{}

func initCommands() {
	commands["help"] = cliCommand{
		name:        "help",
		description: "Display a heko message",
		callback:    commandHelp,
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	commands["map"] = cliCommand{
		name:        "map",
		description: "Shows the next 20 location areas",
		callback:    commandMap,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Shows the previous 20 location areas",
		callback:    commandMapb,
	}
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Shows the pokemon in the specified area",
		callback:    commandExplore,
	}
	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Tries to catch the specified pokemon",
		callback:    commandCatch,
	}
	commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "Shows details of specified pokemon if caught",
		callback:    commandInspect,
	}
	commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "Displays all pokemon that you caught so far",
		callback:    commandPokedex,
	}
}

func commandExit(cfg *config.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config.Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, v := range commands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(cfg *config.Config) error {
	result, err := pokeapi.GetNextLocationAreaBatch(cfg)
	if err != nil {
		return err
	}

	for _, location := range result {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(cfg *config.Config) error {
	result, err := pokeapi.GetPreviousLocationAreaBatch(cfg)
	if err != nil {
		return err
	}

	for _, location := range result {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(cfg *config.Config) error {
	fmt.Printf("Exploring %s...", "")
	result, err := pokeapi.GetPokemonInLocation(cfg)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, v := range result.PokemonEncounters {
		fmt.Printf("- %s\n", v.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config.Config) error {
	pokemonName := cfg.Input[1]
	if pokemonName == "" {
		return errors.New("no pokemon specified")
	}
	pokemon, err := pokeapi.CatchPokemon(pokemonName)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	random := rand.IntN(pokemon.BaseExperience)
	if float64(random) >= float64(pokemon.BaseExperience)*0.7 {
		fmt.Printf("%s was caught\n", pokemonName)
		pokedex.Pokedex[pokemonName] = pokemon
		return nil
	}

	fmt.Printf("%s escaped!\n", pokemonName)
	return nil
}

func commandInspect(cfg *config.Config) error {
	pokemonName := cfg.Input[1]
	if pokemonName == "" {
		return errors.New("No pokemon specified")
	}

	if pokemon, ok := pokedex.Pokedex[pokemonName]; ok {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range pokemon.Stats {
			fmt.Printf("-%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, t := range pokemon.Types {
			fmt.Printf("- %s\n", t.Type.Name)
		}
		return nil
	}

	fmt.Println("You did not caught that pokemon yet")

	return nil
}

func commandPokedex(cfg *config.Config) error {
	for _, pokemon := range pokedex.Pokedex {
		fmt.Println(pokemon.Name)
	}
	return nil
}
