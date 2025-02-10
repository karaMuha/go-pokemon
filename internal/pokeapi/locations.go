package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/karaMuha/pokedex/internal/config"
	"github.com/karaMuha/pokedex/internal/pokecache"
)

type LocationArea struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetNextLocationAreaBatch(cfg *config.Config) ([]Location, error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.Next != "" {
		url = cfg.Next
	}

	if data, ok := cfg.Cache.Get(url); ok {
		var responseData LocationArea
		err := json.Unmarshal(data, &responseData)
		if err != nil {
			return []Location{}, err
		}

		return responseData.Results, nil
	}

	locationArea, err := callAPI(url, cfg.Cache)
	if err != nil {
		return []Location{}, err
	}

	cfg.Next = locationArea.Next
	cfg.Previous = locationArea.Previous
	return locationArea.Results, nil
}

func GetPreviousLocationAreaBatch(cfg *config.Config) ([]Location, error) {
	if cfg.Previous == "" {
		return []Location{}, errors.New("You are at the start, cannot go back")
	}

	url := cfg.Previous
	locationArea, err := callAPI(url, cfg.Cache)
	if err != nil {
		return []Location{}, err
	}

	cfg.Next = locationArea.Next
	cfg.Previous = locationArea.Previous
	return locationArea.Results, nil
}

func callAPI(url string, cache *pokecache.Cache) (LocationArea, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationArea{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, err
	}

	cache.Add(url, body)

	var responseData LocationArea
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return LocationArea{}, err
	}

	return responseData, nil
}
