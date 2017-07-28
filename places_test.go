package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetPlacesSuccess(t *testing.T) {
	location := Location{
		Latitude:  37.483872693672,
		Longitude: -122.14900441942,
	}

	expected := []Place{Place{Name: "Bar Marsella"}}

	googleServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(newGooglePlacesSearchResponse(expected))
		}),
	)
	defer googleServer.Close()
	client := GooglePlacesClient{
		BaseURL: googleServer.URL,
	}
	got, err := location.GetPlaces(client)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestGetPlaceDetailsSuccess(t *testing.T) {
	place := Place{
		ID:   "stn46SGNR452sfg",
		Name: "BarMarsella",
	}

	expected := Place{
		ID:      place.ID,
		Name:    place.Name,
		Rating:  4,
		Website: "www.example.com",
		Location: Location{
			Latitude:  37.483872693672,
			Longitude: -122.14900441942,
		},
	}

	googleServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(newGooglePlaceDetailsResponse(expected))
		}),
	)
	defer googleServer.Close()
	client := GooglePlacesClient{
		BaseURL: googleServer.URL,
	}

	err := place.GetDetails(client)
	if err != nil {
		t.Errorf("unexpected error getting details:", err)
	}

	if !reflect.DeepEqual(place, expected) {
		t.Errorf("expected %v, got %v", expected, place)
	}
}

func newGooglePlacesSearchResponse(places []Place) GooglePlacesSearchResponse {
	return GooglePlacesSearchResponse{
		Results: []Place{
			Place{
				Name: places[0].Name,
			},
		},
	}
}

func newGooglePlaceDetailsResponse(place Place) GooglePlacesDetailsResponse {
	return GooglePlacesDetailsResponse{
		Place: Place{
			Website: place.Website,
			Rating:  place.Rating,
			Geometry: Geometry{
				Location: place.Location,
			},
		},
	}
}

func mustGetResponse(t *testing.T, url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	return resp
}
