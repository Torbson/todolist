package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

// Router
func muxRouter() http.Handler {
	// Create Router
	r := mux.NewRouter()
	r.HandleFunc("/", getIndex).Methods("GET")
	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", deleteTrash).Methods("DELETE")
	r.HandleFunc("/todos", postTodo).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", getTodo).Methods("GET")
	r.HandleFunc("/todos/{id:[0-9]+}", putTodo).Methods("PUT")
	r.HandleFunc("/todos/{id:[0-9]+}", deleteTodo).Methods("DELETE")
	return r
}

// ROUTES
func getIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todos", http.StatusMovedPermanently)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	// Get todos batch
	var todos []Todo
	db.Preload("Tasks").Find(&todos)

	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(todos)
}

func postTodo(w http.ResponseWriter, r *http.Request) {
	// Get Input
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	// validate JSON input
	if errs := validator.Validate(todo); errs != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(errs)
		return
	}
	// Post todo
	db.Create(&todo)
	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(todo)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	// Get Input
	params := mux.Vars(r)
	// Check input
	id, err := strconv.ParseUint(params["id"], 10, 0)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	// Get todo by id
	var todo Todo
	db.Preload("Tasks").First(&todo, id)
	// Check todo
	if 0 == todo.ID {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(todo)
}

func putTodo(w http.ResponseWriter, r *http.Request) {
	// Get input
	params := mux.Vars(r)
	// Check input
	id, err := strconv.ParseUint(params["id"], 10, 0)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var todo Todo
	db.Preload("Tasks").First(&todo, id)
	// Check todo
	if 0 == todo.ID {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get JSON input
	var todo_update Todo
	json.NewDecoder(r.Body).Decode(&todo_update)
	// validate JSON input
	if errs := validator.Validate(todo_update); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}
	// Check if JSON matches with path
	if todo_update.ID != uint(id) {
		fmt.Println("resource path mismatch JSON id!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Put todo
	db.Save(&todo_update)

	// get updated todo
	db.Preload("Tasks").First(&todo, id)

	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	// Get input
	params := mux.Vars(r)
	// Check input
	id, err := strconv.ParseUint(params["id"], 10, 0)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Get todo by id
	var todo Todo
	db.Preload("Tasks").First(&todo, id)
	// Check todo
	if 0 == todo.ID {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Move todo to trash (soft delete)
	//db.Select("Tasks").Delete(&todo, params["id"]) // soft-delete cascade Tasks
	db.Delete(&todo, params["id"]) // soft-delete
	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(todo)
}

func deleteTrash(w http.ResponseWriter, r *http.Request) {
	// Get trash (soft delete) todos
	var todo []Todo
	db.Unscoped().Where("deleted_at is NOT NULL").Preload("Tasks").Find(&todo)
	// empty trash (permanent delete) todos
	db.Unscoped().Select("Tasks").Delete(&todo)
	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(todo)
}
