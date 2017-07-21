package main

import (
	"fmt"
	"os"
)

const ErrInvalidCoordinates = "invalid coordinates"

type GooglePlacesClient struct {
	BaseURL string
}

func NewGooglePlacesClient(c Config) GooglePlacesClient {
	client := GooglePlacesClient{
		BaseURL: fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=-33.8670,151.1957&radius=500&types=food&key=%s", os.Getenv("GOOGLE_API_KEY")),
	}
	if c.APIBaseURL != "" {
		client.BaseURL = c.APIBaseURL
	}
	return client
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
