package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type TestStruct struct {
	name   string
	input  TestInput
	output TestOutput
}

type TestInput struct {
	method      string
	path        string
	contentType string
	content     io.Reader
}
type TestOutput struct {
	status      int
	contentType string
	content     string
}

func Test_Router(t *testing.T) {
	tests := []TestStruct{
		{name: "GET /", input: TestInput{method: http.MethodGet, path: "/", content: nil}, output: TestOutput{status: http.StatusOK}},
		{name: "GET /todos", input: TestInput{method: http.MethodGet, path: "/todos", content: nil}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos WITHOUT JSON", input: TestInput{method: http.MethodPost, path: "/todos", content: nil}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "POST /todos", input: TestInput{method: http.MethodPost, path: "/todos", content: nil}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos", input: TestInput{method: http.MethodPost, path: "/todos", content: nil}, output: TestOutput{status: http.StatusOK}},
	}
	// get environment
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Connect to database
	psql_user := os.Getenv("POSTGRES_USER")
	psql_pw := os.Getenv("POSTGRES_PASSWORD")
	psql_db := os.Getenv("POSTGRES_DB")
	psql_host := os.Getenv("POSTGRES_HOST")
	psql_port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s_test port=%s", psql_host, psql_user, psql_pw, psql_db, psql_port)
	connect_db(dsn)
	connect_db("host=localhost user=postgres password=PostgresTest dbname=todolist_test port=5432")
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	for _, test := range tests {
		//create request based on test table entry
		req, err := http.NewRequest(test.input.method, fmt.Sprintf("%s%s", srv.URL, test.input.path), test.input.content)
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
			body, err := ioutil.ReadAll(r.Body)
			content := string(body)
			if err != nil {
				t.Fatalf("get response body error: %v", err)
			}
			// check content
			if test.output.content != content {
				t.Errorf("expected content: %v; got %v", test.output.content, content)
			}
		}
		t.Logf("%s SUCCESS", test.name)
	}
}
