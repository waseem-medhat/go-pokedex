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
