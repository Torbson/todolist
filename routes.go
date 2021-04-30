package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
	// Post todo
	db.Create(&todo)
	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(todo)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	// Get Input
	params := mux.Vars(r)
	// Get todo by id
	var todo Todo
	db.Preload("Tasks").First(&todo, params["id"])
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
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	// Check Input
	id, err := strconv.ParseUint(params["id"], 10, 0)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if todo.ID != uint(id) {
		fmt.Println("ID of JSON does not match with resource path!")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	// Put todo
	db.Save(&todo)
	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	// Get Input
	params := mux.Vars(r)
	// Get todo by id
	var todo Todo
	db.Preload("Tasks").First(&todo, params["id"])
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
