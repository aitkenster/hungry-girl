package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	APIBaseURL string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	latitude := 37.483872693672
	longitude := -122.14900441942
	location := Location{latitude, longitude}
	client := NewGooglePlacesClient(Config{})
	places, err := location.GetPlaces(client)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(places)
}