package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	ID      string
	Name    string
	Website string
	Rating  float64
}

type GooglePlacesClient struct {
	BaseURL string
	APIKey  string
}

type errorResponse struct {
	Message string `json:"message"`
}

type GooglePlacesSearchResponse struct {
	Results []Result `json:"results"`
}

type GooglePlacesDetailsResponse struct {
	Result Result `json:"result"`
}

type Result struct {
	Name    string  `json:"name"`
	ID      string  `json:"place_id"`
	Website string  `json:"website"`
	Rating  float64 `json:"rating"`
}

func (l Location) GetPlaces(client GooglePlacesClient) ([]Place, error) {
	url := fmt.Sprintf("%s/nearbysearch/json?location=%v,%v&radius=500&types=food&key=%s", client.BaseURL, l.Latitude, l.Longitude, client.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return nil, errors.New("error retrieving Google Places response")
	}
	var g GooglePlacesSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		fmt.Println(err)
	}
	var p []Place
	numPlaces := 0
	for _, result := range g.Results {
		if numPlaces < placesLimit {
			p = append(p, Place{
				Name: result.Name,
				ID:   result.ID,
			})
			numPlaces++
		}
	}
	return p, nil
}

func (p *Place) GetWebsite(client GooglePlacesClient) error {
	url := fmt.Sprintf("%s/details/json?placeid=%s&key=%s", client.BaseURL, p.ID, client.APIKey)
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return errors.New("error retrieving Google Places response")
	}
	var g GooglePlacesDetailsResponse
	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		fmt.Println(err)
		return err
	}
	p.Website = g.Result.Website
	p.Rating = g.Result.Rating
	return nil
}

func NewGooglePlacesClient(c Config) GooglePlacesClient {
	client := GooglePlacesClient{
		BaseURL: "https://maps.googleapis.com/maps/api/place",
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
