package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/karaMuha/pokedex/internal/config"
	"github.com/karaMuha/pokedex/internal/models"
	"github.com/karaMuha/pokedex/internal/pokecache"
)

func GetNextLocationAreaBatch(cfg *config.Config) ([]models.Location, error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.Next != "" {
		url = cfg.Next
	}

	if data, ok := cfg.Cache.Get(url); ok {
		var responseData models.LocationArea
		err := json.Unmarshal(data, &responseData)
		if err != nil {
			return []models.Location{}, err
		}

		return responseData.Results, nil
	}

	locationArea, err := callAPI(url, cfg.Cache)
	if err != nil {
		return []models.Location{}, err
	}

	cfg.Next = locationArea.Next
	cfg.Previous = locationArea.Previous
	return locationArea.Results, nil
}

func GetPreviousLocationAreaBatch(cfg *config.Config) ([]models.Location, error) {
	if cfg.Previous == "" {
		return []models.Location{}, errors.New("You are at the start, cannot go back")
	}

	url := cfg.Previous
	locationArea, err := callAPI(url, cfg.Cache)
	if err != nil {
		return []models.Location{}, err
	}

	cfg.Next = locationArea.Next
	cfg.Previous = locationArea.Previous
	return locationArea.Results, nil
}

func callAPI(url string, cache *pokecache.Cache) (models.LocationArea, error) {
	res, err := http.Get(url)
	if err != nil {
		return models.LocationArea{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.LocationArea{}, err
	}

	cache.Add(url, body)

	var responseData models.LocationArea
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return models.LocationArea{}, err
	}

	return responseData, nil
}
