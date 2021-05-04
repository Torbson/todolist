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
	r.Handle("/", requires_api_key(http.HandlerFunc(getTodos))).Methods("GET")
	r.Handle("/todos", requires_api_key(http.HandlerFunc(getTodos))).Methods("GET")
	r.Handle("/todos", requires_api_key(http.HandlerFunc(postTodo))).Methods("POST")
	r.Handle("/todos", requires_api_key(http.HandlerFunc(deleteTrash))).Methods("DELETE")
	r.Handle("/todos/{id:[0-9]+}", requires_api_key(http.HandlerFunc(getTodo))).Methods("GET")
	r.Handle("/todos/{id:[0-9]+}", requires_api_key(http.HandlerFunc(putTodo))).Methods("PUT")
	r.Handle("/todos/{id:[0-9]+}", requires_api_key(http.HandlerFunc(deleteTodo))).Methods("DELETE")
	return r
}

// ROUTES
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

	// sanitize todo strings for html / javascript
	html_sanitize_todo(&todo)

	// validate JSON input
	if errs := validator.Validate(todo); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(errs)
		return
	}

	//check if ids are set in todo
	if 0 != todo.ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check if ids are set in todo.Tasks
	if nil != todo.Tasks {
		for _, task := range todo.Tasks {
			if 0 != task.ID {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
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

	// sanitize todo strings for html / javascript
	html_sanitize_todo(&todo_update)
	// validate JSON input
	if errs := validator.Validate(todo_update); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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

	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// get updated todo
	db.Preload("Tasks").First(&todo, id)
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
	db.Delete(&todo, params["id"]) // soft-delete
	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(todo)
}

func deleteTrash(w http.ResponseWriter, r *http.Request) {
	// Get trash (soft delete) todos
	var todos []Todo
	db.Unscoped().Where("deleted_at is NOT NULL").Preload("Tasks").Find(&todos)
	// empty trash (permanent delete) todos
	db.Unscoped().Select("Tasks").Delete(&todos)
	// Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(todos)
}
