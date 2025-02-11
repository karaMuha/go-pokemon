package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/karaMuha/pokedex/internal/config"
	"github.com/karaMuha/pokedex/internal/pokecache"
)

func main() {
	initCommands()
	input := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5 * time.Second)
	cfg := config.NewConfig(cache)

	for {
		fmt.Print("Pokedex > ")
		input.Scan()
		if err := input.Err(); err != nil {
			fmt.Println("Something went wrong")
			os.Exit(1)
		}
		text := input.Text()
		textLower := strings.ToLower(text)
		clearedText := strings.Fields(textLower)
		command, ok := commands[clearedText[0]]
		if !ok {
			fmt.Println("Unknown command")
		}
		cfg.Input = clearedText
		if err := command.callback(cfg); err != nil {
			fmt.Println(err)
		}
	}
}
