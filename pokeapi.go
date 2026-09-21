package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Enri258/pokedexcli/internal/pokecache"
)

const locationAreaBaseURL = "https://pokeapi.co/api/v2/location-area/"
const pokemonBaseURL = "https://pokeapi.co/api/v2/pokemon/"

type locationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type locationAreaResponse struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func getCachedData(
	url string,
	cache *pokecache.Cache,
) ([]byte, error) {
	data, ok := cache.Get(url)
	if ok {
		return data, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected status code: %d",
			res.StatusCode,
		)
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	cache.Add(url, data)

	return data, nil
}

func getLocationAreas(
	url string,
	cache *pokecache.Cache,
) (locationAreasResponse, error) {
	data, err := getCachedData(url, cache)
	if err != nil {
		return locationAreasResponse{}, err
	}

	var locations locationAreasResponse

	err = json.Unmarshal(data, &locations)
	if err != nil {
		return locationAreasResponse{}, err
	}

	return locations, nil
}

func getLocationArea(
	name string,
	cache *pokecache.Cache,
) (locationAreaResponse, error) {
	url := locationAreaBaseURL + name

	data, err := getCachedData(url, cache)
	if err != nil {
		return locationAreaResponse{}, err
	}

	var area locationAreaResponse

	err = json.Unmarshal(data, &area)
	if err != nil {
		return locationAreaResponse{}, err
	}

	return area, nil
}

func getPokemon(
	name string,
	cache *pokecache.Cache,
) (Pokemon, error) {
	url := pokemonBaseURL + name

	data, err := getCachedData(url, cache)
	if err != nil {
		return Pokemon{}, err
	}

	var pokemon Pokemon

	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
