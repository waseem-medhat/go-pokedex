package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationData struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type ExploreData struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func GetLocations(url string) (LocationData, error) {
	var locations LocationData

	r, err := http.Get(url)
	if err != nil {
		return locations, err
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return locations, err
	}

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return locations, err
	}

	return locations, err
}

func GetExplore(url string) (ExploreData, error) {
	var exploreData ExploreData

	r, err := http.Get(url)
	if err != nil {
		return exploreData, err
	}

	if r.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("Not found")
		return exploreData, err
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return exploreData, err
	}

	err = json.Unmarshal(body, &exploreData)
	if err != nil {
		return exploreData, err
	}

	return exploreData, err
}

func GetPokemon(url string) (Pokemon, error) {
	var pokemon Pokemon

	r, err := http.Get(url)
	if err != nil {
		return pokemon, err
	}

	if r.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("Not found")
		return pokemon, err
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return pokemon, err
	}

	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return pokemon, err
	}

	return pokemon, err
}
