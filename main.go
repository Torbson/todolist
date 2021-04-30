package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	//"math/rand"
	//"strconv"
	"github.com/gorilla/mux"
)

// Gorm Database
const dsn = "host=localhost user=postgres password=PostgresTest dbname=todolist port=5432"

var db *gorm.DB
var err error

// MAIN
func main() {
	// Init PostgreSQL
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	// Auto migrate in case of schema / model changes
	db.AutoMigrate(&Todo{}, &Task{})

	// Init Router
	r := mux.NewRouter()

	// Route Handlers
	r.HandleFunc("/", getIndex).Methods("GET")
	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", deleteTrash).Methods("DELETE")
	r.HandleFunc("/todos", postTodo).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", getTodo).Methods("GET")
	r.HandleFunc("/todos/{id:[0-9]+}", putTodo).Methods("PUT")
	r.HandleFunc("/todos/{id:[0-9]+}", deleteTodo).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}
