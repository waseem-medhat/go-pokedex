package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type LocationData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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
