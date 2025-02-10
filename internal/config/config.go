package config

import "github.com/karaMuha/pokedex/internal/pokecache"

type Config struct {
	Input        []string
	Next         string
	Previous     string
	Cache        *pokecache.Cache
	LocationArea string
}

func NewConfig(cache *pokecache.Cache) *Config {
	return &Config{
		Cache: cache,
	}
}
