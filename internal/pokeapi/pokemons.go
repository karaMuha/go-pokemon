package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/karaMuha/pokedex/internal/config"
	"github.com/karaMuha/pokedex/internal/models"
)

func GetPokemonInLocation(cfg *config.Config) (models.PokemonInLocation, error) {
	location := cfg.Input[1]
	if location == "" {
		return models.PokemonInLocation{}, errors.New("No location specified")
	}
	url := "https://pokeapi.co/api/v2/location-area/" + location

	res, err := http.Get(url)
	if err != nil {
		return models.PokemonInLocation{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.PokemonInLocation{}, err
	}

	cfg.Cache.Add(url, body)

	var responseData models.PokemonInLocation
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return models.PokemonInLocation{}, err
	}

	return responseData, nil
}

func CatchPokemon(pokemonName string) (models.Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	res, err := http.Get(url)
	if err != nil {
		return models.Pokemon{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Pokemon{}, err
	}

	var pokemon models.Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return models.Pokemon{}, err
	}

	return pokemon, nil
}
