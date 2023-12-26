package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type locationData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocations(url string) (locationData, error) {
    var locations locationData

	r, err := http.Get(url)
	if err != nil {
		return locations, err
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return locations, err
	}

	json.Unmarshal(body, &locations)
    return locations, nil
}
