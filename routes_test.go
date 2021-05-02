package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type Todos struct {
	Todos []Todo
}

type TestStruct struct {
	name   string
	input  TestInput
	output TestOutput
}

type TestInput struct {
	method      string
	path        string
	contentType string
	content     string
}
type TestOutput struct {
	status      int
	contentType string
	content     string
}

func Test_Router(t *testing.T) {
	tests := []TestStruct{
		{name: "GET /", input: TestInput{method: http.MethodGet, path: "/"}, output: TestOutput{status: http.StatusOK}},
		{name: "GET /todos", input: TestInput{method: http.MethodGet, path: "/todos"}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos", input: TestInput{method: http.MethodPost, path: "/todos", content: json_todo1_post}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos only name", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"name" : "minimal requirements"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos without JSON", input: TestInput{method: http.MethodPost, path: "/todos"}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "POST /todos without name", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"description" : "name required"}`}, output: TestOutput{status: http.StatusBadRequest}},
	}
	/*
		r.HandleFunc("/", getIndex).Methods("GET")
		r.HandleFunc("/todos", getTodos).Methods("GET")
		r.HandleFunc("/todos", postTodo).Methods("POST")
		r.HandleFunc("/todos", deleteTrash).Methods("DELETE")
		r.HandleFunc("/todos/{id:[0-9]+}", getTodo).Methods("GET")
		r.HandleFunc("/todos/{id:[0-9]+}", putTodo).Methods("PUT")
		r.HandleFunc("/todos/{id:[0-9]+}", deleteTodo).Methods("DELETE")
	*/
	// check env
	check_env()
	// get env
	psql_user := os.Getenv("POSTGRES_USER")
	psql_pw := os.Getenv("POSTGRES_PASSWORD")
	psql_db := os.Getenv("POSTGRES_DB")
	psql_host := os.Getenv("POSTGRES_HOST")
	psql_port := os.Getenv("POSTGRES_PORT")
	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s_test port=%s", psql_host, psql_user, psql_pw, psql_db, psql_port)
	connect_db(dsn)
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	for i, test := range tests {
		var body io.Reader
		body = nil
		if 0 < len(test.input.content) {
			body = strings.NewReader(test.input.content)
		}
		//create request based on test table entry
		req, err := http.NewRequest(test.input.method, fmt.Sprintf("%s%s", srv.URL, test.input.path), body)
		if err != nil {
			t.Fatalf("create request error: %v", err)
		}
		// send request
		r, err := client.Do(req)
		if err != nil {
			t.Fatalf("send request error: %v", err)
		}
		defer r.Body.Close()

		// status code expected?
		if 100 <= test.output.status {
			// check status code
			if test.output.status != r.StatusCode {
				t.Errorf("expected status code: %v; got %v from %s%s", test.output.status, r.StatusCode, srv.URL, test.input.path)
			}
		}

		// content expected?
		if 0 < len(test.output.content) {
			// get content
			content, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("get response body error: %v", err)
			}
			// get expected content
			expected := new(bytes.Buffer)
			err = json.Compact(expected, []byte(test.output.content))
			if err != nil {
				fmt.Println(err)
			}
			// check content
			if !bytes.Contains(content, expected.Bytes()) {
				t.Errorf("expected content: %v; got %v", content, expected.Bytes())
			}
		}
		t.Logf("Test_Router tabbletest: %d, name: %s", i, test.name)
	}
}

var json_todo1_post string = `
{
	"name" : "todolist go REST API",
	"description" : "write a REST API in go for a simple todo list. Each todo can have multiple subtodos called tasks",	
    "tasks":[
        {
            "name" : "Create Server",
            "description" : "Create Server according to specs with go"
        },
        {
            "name" : "Unittest",
            "description" : "Build Unittest for go server"
        },
        {
            "name" : "Docker",
            "description" : "create docker container with go server"
        },
        {
            "name" : "Postman",
            "description" : "create postman test for docker container"
        },
        {
            "name" : "Documentation",
            "description" : "write markdown documentation that describes installation, API, models, dependencies, deployment and so on"
        }
    ]
}`
