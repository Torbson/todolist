package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// ENV vars:
var PORT string
var POSTGRES_USER string
var POSTGRES_PASSWORD string
var POSTGRES_DB string
var POSTGRES_HOST string
var POSTGRES_PORT string
var TODOLIST_API_KEY string

func get_env() {
	check_env()
	PORT = os.Getenv("PORT")
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB = os.Getenv("POSTGRES_DB")
	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
	TODOLIST_API_KEY = os.Getenv("TODOLIST_API_KEY")
}

func check_env() {
	// check env
	env := os.Getenv("ENV")
	// if env variable is not set load environment from .env file
	if "there" != env {
		// load env
		log.Print("load env")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

// MAIN
func main() {
	// get env
	get_env()

	// Connect to database
	connect_db(POSTGRES_DB)

	// Start http server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), muxRouter()))
}
