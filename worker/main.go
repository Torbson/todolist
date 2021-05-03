package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func connect_db(dsn string) {
	// Connect PostgreSQL
	log.Print("connect to Postgres ...")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	log.Print("successfully connected to Postgres")
	// Auto migrate in case of schema / model changes
	db.AutoMigrate(&Todo{}, &Task{})
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
	// check env
	check_env()
	// get env
	port := os.Getenv("PORT")
	psql_user := os.Getenv("POSTGRES_USER")
	psql_pw := os.Getenv("POSTGRES_PASSWORD")
	psql_db := os.Getenv("POSTGRES_DB")
	psql_host := os.Getenv("POSTGRES_HOST")
	psql_port := os.Getenv("POSTGRES_PORT")

	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", psql_host, psql_user, psql_pw, psql_db, psql_port)
	connect_db(dsn)

	// Start http server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), muxRouter()))
}
