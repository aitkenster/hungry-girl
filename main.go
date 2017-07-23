package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
