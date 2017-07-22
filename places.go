package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

const (
	placesLimit           = 3
	ErrInvalidCoordinates = "invalid coordinates"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type Place struct {
	Name    string
	Stars   int
	Website string
}

type GooglePlacesClient struct {
	BaseURL string
	APIKey  string
}

type errorResponse struct {
	Message string `json:"message"`
}

type GooglePlacesResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Name string `json:"name"`
}

func (l Location) GetPlaces() ([]Place, error) {
	client := NewGooglePlacesClient(Config{})
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
	var p []Place
	numPlaces := 0
	for _, result := range g.Results {
		if numPlaces < placesLimit {
			p = append(p, Place{Name: result.Name})
			numPlaces++
		}
	}
	return p, nil
}

func NewGooglePlacesClient(c Config) GooglePlacesClient {
	client := GooglePlacesClient{
		BaseURL: "https://maps.googleapis.com/maps/api/place/nearbysearch/json?",
		APIKey:  os.Getenv("GOOGLE_API_KEY"),
	}
	if c.APIBaseURL != "" {
		client.BaseURL = c.APIBaseURL
	}
	return client
}

func NewLocation(latitude, longitude float64) (*Location, error) {
	l := Location{
		Latitude:  latitude,
		Longitude: longitude,
	}

	return &l, nil
}
