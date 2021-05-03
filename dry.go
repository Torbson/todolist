package main

//Helper functions

import (
	"log"
	"net/http"

	"github.com/microcosm-cc/bluemonday"
)

var sanitizer *bluemonday.Policy

func html_sanitize_todo(todo *Todo) {
	if nil == sanitizer {
		sanitizer = bluemonday.UGCPolicy()
	}
	// sanitize todo strings for html / javascript
	todo.Name, todo.Description = sanitizer.Sanitize(todo.Name), sanitizer.Sanitize(todo.Description)
	for _, task := range todo.Tasks {
		task.Name, task.Description = sanitizer.Sanitize(task.Name), sanitizer.Sanitize(task.Description)
	}
}

// decorator for simple api key auth / identification
func requires_api_key(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remote_ip := getIP(r)
		keys, ok := r.URL.Query()["key"]
		//check if key is set as query parameter
		if !ok || len(keys[0]) < 1 {
			log.Printf("WARNING: IP %s tried to access critical API without API key!", remote_ip)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//check if key matches (can be extended to search match in multiple client api keys)
		if keys[0] != TODOLIST_API_KEY {
			log.Printf("WARNING: IP %s tried to access critical API with wrong API key!", remote_ip)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
