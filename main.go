package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Gorm Database
const dsn = "host=localhost user=postgres password=PostgresTest dbname=todolist port=5432"

var db *gorm.DB
var err error

func init_db(dsn string) {
	// Init PostgreSQL
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	// Auto migrate in case of schema / model changes
	db.AutoMigrate(&Todo{}, &Task{})
}

// MAIN
func main() {
	init_db(dsn)

	// Start server
	log.Fatal(http.ListenAndServe(":8000", muxRouter()))
}
