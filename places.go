package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
	resp, err := http.Get(client.BaseURL)
	if err != nil {
		fmt.Println(err)
	}
	var g GooglePlacesResponse
	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		fmt.Println(err)
	}
	p := Placelist{
		Place{
			Name: g.Results[0].Name,
		},
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
