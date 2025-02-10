package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/karaMuha/pokedex/internal/config"
)

type PokemonInLocation struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []Stat `json:"stats"`
	Types          []Type `json:"types"`
}

type Species struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Stat struct {
	BaseStat int     `json:"base_stat"`
	Effort   int     `json:"effort"`
	Stat     Species `json:"stat"`
}

type Type struct {
	Slot int     `json:"slot"`
	Type Species `json:"type"`
}

func GetPokemonInLocation(cfg *config.Config) (PokemonInLocation, error) {
	location := cfg.Input[1]
	if location == "" {
		return PokemonInLocation{}, errors.New("No location specified")
	}
	url := "https://pokeapi.co/api/v2/location-area/" + location

	res, err := http.Get(url)
	if err != nil {
		return PokemonInLocation{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return PokemonInLocation{}, err
	}

	cfg.Cache.Add(url, body)

	var responseData PokemonInLocation
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return PokemonInLocation{}, err
	}

	return responseData, nil
}

func CatchPokemon(pokemonName string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
