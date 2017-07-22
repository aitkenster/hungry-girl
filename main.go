package main

import (
	"fmt"
	"log"
	"net/http"

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
	http.HandleFunc("/messenger", MessengerRequestHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
