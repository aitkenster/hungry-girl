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
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Place struct {
	Name     string   `json:"name"`
	ID       string   `json:"place_id"`
	Website  string   `json:"website"`
	Rating   float64  `json:"rating"`
	Geometry Geometry `json:"geometry"`
	Location Location
}

type GooglePlacesClient struct {
	BaseURL string
	APIKey  string
}

type errorResponse struct {
	Message string `json:"message"`
}

type GooglePlacesSearchResponse struct {
	Results []Place `json:"results"`
}

type GooglePlacesDetailsResponse struct {
	Place Place `json:"result"`
}

type Geometry struct {
	Location Location `json:"location"`
}

func (l Location) GetPlacesFromGoogle(client GooglePlacesClient) ([]Place, error) {
	url := fmt.Sprintf("%s/nearbysearch/json?location=%v,%v&radius=500&type=restaurant&opennow=true&key=%s", client.BaseURL, l.Latitude, l.Longitude, client.APIKey)

	resp, err := getSuccessfulResponseFromGooglePlaces(url)
	if err != nil {
		return nil, err
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

func (p *Place) GetDetails(client GooglePlacesClient) error {
	url := fmt.Sprintf("%s/details/json?placeid=%s&key=%s", client.BaseURL, p.ID, client.APIKey)
	resp, err := getSuccessfulResponseFromGooglePlaces(url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var g GooglePlacesDetailsResponse
	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		fmt.Println(err)
		return err
	}

	p.Website = g.Place.Website
	p.Rating = g.Place.Rating
	p.Location = g.Place.Geometry.Location
	return nil
}

func NewGooglePlacesClient(c Config) GooglePlacesClient {
	client := GooglePlacesClient{
		BaseURL: "https://maps.googleapis.com/maps/api/place",
		APIKey:  os.Getenv("GOOGLE_PLACES_API_KEY"),
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

func getSuccessfulResponseFromGooglePlaces(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error retrieving Google Places response")
	}
	return resp, nil
}

func (p *Place) StaticMapUrl() string {
	return fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?markers=color:red|label:B|%v,%v&size=360x360&zoom=13", p.Location.Latitude, p.Location.Longitude)
}

func (p *Place) LinkMapUrl() string {
	return fmt.Sprintf("https://www.google.com/maps/place/?q=place_id:%s", p.ID)
}
