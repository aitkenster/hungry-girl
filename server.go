package main

import "os"

const ErrInvalidCoordinates = "invalid coordinates"

type GooglePlacesClient struct {
	BaseURL string
	APIKey  string
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

type errorResponse struct {
	Message string `json:"message"`
}

type GooglePlacesResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Name string `json:"name"`
}
