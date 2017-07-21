package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const placesLimit = 3

type Location struct {
	Latitude  float64
	Longitude float64
}

type Placelist []Place

type Place struct {
	Name    string
	Stars   int
	Website string
}

func (l Location) GetPlaces(client GooglePlacesClient) (*Placelist, error) {
	url := fmt.Sprintf("%slocation=%v,%v&radius=500&types=food&key=%s", client.BaseURL, l.Latitude, l.Longitude, client.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return nil, errors.New("error retrieving Google Places response")
	}
	var g GooglePlacesResponse
	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		fmt.Println(err)
	}
	var p Placelist
	numPlaces := 0
	for _, result := range g.Results {
		if numPlaces < placesLimit {
			p = append(p, Place{Name: result.Name})
			numPlaces++
		}
	}
	return &p, nil
}

func NewLocation(latitude, longitude float64) (*Location, error) {
	l := Location{
		Latitude:  latitude,
		Longitude: longitude,
	}

	return &l, nil
}
