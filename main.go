package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	APIBaseURL string
}

var DB *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	http.HandleFunc("/messenger", MessengerRequestHandler)

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
