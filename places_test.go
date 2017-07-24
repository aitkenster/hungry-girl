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

func newGooglePlacesSearchResponse(places []Place) GooglePlacesSearchResponse {
	return GooglePlacesSearchResponse{
		Results: []Result{
			Result{
				Name: places[0].Name,
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
