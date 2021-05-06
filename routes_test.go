package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

var testCount int = 1

var json_todo_post string = `
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

var json_todo_1941041_post string = `
{
	"id" : 1941041,
	"name" : "todolist go REST API",
	"description" : "write a REST API in go for a simple todo list. Each todo can have multiple subtodos called tasks",	
    "tasks":[
        {
			"id" : 1941041,
            "name" : "Create Server",
            "description" : "Create Server according to specs with go"
        },
        {
			"id" : 1941042,
            "name" : "Unittest",
            "description" : "Build Unittest for go server"
        },
        {
			"id" : 1941043,
            "name" : "Docker",
            "description" : "create docker container with go server"
        },
        {
			"id" : 1941044,
            "name" : "Postman",
            "description" : "create postman test for docker container"
        },
        {
			"id" : 1941045,
            "name" : "Documentation",
            "description" : "write markdown documentation that describes installation, API, models, dependencies, deployment and so on"
        }
    ]
}`

func init_test_db() {
	connect_db(fmt.Sprintf("%s_test", POSTGRES_DB))
}

func Test_Router(t *testing.T) {
	tests := []TestStruct{
		{name: "GET /", input: TestInput{method: http.MethodGet, path: "/"}, output: TestOutput{status: http.StatusOK}},
		{name: "GET /todos", input: TestInput{method: http.MethodGet, path: "/todos"}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (add todo ignore unknown key)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"name":"check upwards compatibility","unknown unknowns":"ignore until update"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (add todo with name only)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"name":"check post with name only"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (add todo with name empty)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"name":""}`}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "POST /todos (add todo with long name)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"name":"pretty pretty loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooonnnng name"}`}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "POST /todos (add todo and ignore id string)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"id":"100","name":"check post with name only"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (bad request with id 1941041)", input: TestInput{method: http.MethodPost, path: "/todos", content: json_todo_1941041_post}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "POST /todos (bad request without JSON)", input: TestInput{method: http.MethodPost, path: "/todos"}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "POST /todos (bad request without name)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"description":"name required"}`}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "PUT /todos (PUT not supported)", input: TestInput{method: http.MethodPut, path: "/todos", content: `{"name":"check post with name only"}`}, output: TestOutput{status: http.StatusMethodNotAllowed}},
		{name: "PUT /todos/1941042 (PUT not found)", input: TestInput{method: http.MethodPut, path: "/todos/1941042", content: `{"id":1941042,"name":"check post with name only"}`}, output: TestOutput{status: http.StatusNotFound}},
		{name: "PATCH /todos (PATCH not supported)", input: TestInput{method: http.MethodPatch, path: "/todos", content: `{"name":"check post with name only"}`}, output: TestOutput{status: http.StatusMethodNotAllowed}},
		{name: "DELETE /todos (trash empty, ignore JSON)", input: TestInput{method: http.MethodDelete, path: "/todos", content: `{"description":"name required"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "DELETE /todos (trash empty)", input: TestInput{method: http.MethodDelete, path: "/todos"}, output: TestOutput{status: http.StatusOK}},
	}
	// get env
	get_env()
	//init postgres db
	init_test_db()
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	for _, test := range tests {
		var body io.Reader
		body = nil
		if 0 < len(test.input.content) {
			body = strings.NewReader(test.input.content)
		}
		//create request based on test table entry
		req, err := http.NewRequest(test.input.method, fmt.Sprintf("%s%s?key=%s", srv.URL, test.input.path, TODOLIST_API_KEY), body)
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
				t.Errorf("expected status code: %v; got %v", test.output.status, r.StatusCode)
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
			response_json := string(content)
			expected_json := expected.String()
			if response_json != expected_json {
				t.Errorf("expected content: %v;\ngot %v\n", response_json, expected_json)
			}
		}
		t.Logf("Test_Router tabletest: %d, name: %s", testCount, test.name)
		testCount++
	}
}

func Test_Wrong_API_Key(t *testing.T) {
	tests := []TestStruct{
		{name: "GET /", input: TestInput{method: http.MethodGet, path: "/"}, output: TestOutput{status: http.StatusOK}},
		{name: "GET /todos", input: TestInput{method: http.MethodGet, path: "/todos"}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (add todo ignore unknown key)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"name":"check upwards compatibility","unknown unknowns":"ignore until update"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (add todo with name only)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"name":"check post with name only"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (add todo and ignore id string)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"id":"100","name":"check post with name only"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (bad request with id 1941041)", input: TestInput{method: http.MethodPost, path: "/todos", content: json_todo_1941041_post}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "POST /todos (bad request without JSON)", input: TestInput{method: http.MethodPost, path: "/todos"}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "POST /todos (bad request without name)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"description":"name required"}`}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "PUT /todos/1941042 (PUT not found)", input: TestInput{method: http.MethodPut, path: "/todos/1941042", content: `{"id":1941042,"name":"check post with name only"}`}, output: TestOutput{status: http.StatusNotFound}},
		{name: "DELETE /todos (trash empty, ignore JSON)", input: TestInput{method: http.MethodDelete, path: "/todos", content: `{"description":"name required"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "DELETE /todos (trash empty)", input: TestInput{method: http.MethodDelete, path: "/todos"}, output: TestOutput{status: http.StatusOK}},
	}
	// get env
	get_env()
	//init postgres db
	init_test_db()
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	for _, test := range tests {
		var body io.Reader
		body = nil
		if 0 < len(test.input.content) {
			body = strings.NewReader(test.input.content)
		}
		//create request based on test table entry
		req, err := http.NewRequest(test.input.method, fmt.Sprintf("%s%s?key=wrongAPIkey", srv.URL, test.input.path), body)
		if err != nil {
			t.Fatalf("create request error: %v", err)
		}
		// send request
		r, err := client.Do(req)
		if err != nil {
			t.Fatalf("send request error: %v", err)
		}
		defer r.Body.Close()

		// check status code
		if http.StatusUnauthorized != r.StatusCode {
			t.Errorf("expected status code: %v; got %v", http.StatusUnauthorized, r.StatusCode)
		}

		t.Logf("Test_Router tabletest: %d, name: %s", testCount, test.name)
		testCount++
	}
}

func Test_Without_API_Key(t *testing.T) {
	tests := []TestStruct{
		{name: "GET /", input: TestInput{method: http.MethodGet, path: "/"}, output: TestOutput{status: http.StatusOK}},
		{name: "GET /todos", input: TestInput{method: http.MethodGet, path: "/todos"}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (add todo ignore unknown key)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"name":"check upwards compatibility","unknown unknowns":"ignore until update"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (add todo with name only)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"name":"check post with name only"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (add todo and ignore id string)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"id":"100","name":"check post with name only"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "POST /todos (bad request with id 1941041)", input: TestInput{method: http.MethodPost, path: "/todos", content: json_todo_1941041_post}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "POST /todos (bad request without JSON)", input: TestInput{method: http.MethodPost, path: "/todos"}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "POST /todos (bad request without name)", input: TestInput{method: http.MethodPost, path: "/todos", content: `{"description":"name required"}`}, output: TestOutput{status: http.StatusBadRequest}},
		{name: "PUT /todos/1941042 (PUT not found)", input: TestInput{method: http.MethodPut, path: "/todos/1941042", content: `{"id":1941042,"name":"check post with name only"}`}, output: TestOutput{status: http.StatusNotFound}},
		{name: "DELETE /todos (trash empty, ignore JSON)", input: TestInput{method: http.MethodDelete, path: "/todos", content: `{"description":"name required"}`}, output: TestOutput{status: http.StatusOK}},
		{name: "DELETE /todos (trash empty)", input: TestInput{method: http.MethodDelete, path: "/todos"}, output: TestOutput{status: http.StatusOK}},
	}
	// get env
	get_env()
	//init postgres db
	init_test_db()
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	for _, test := range tests {
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

		// check status code
		if http.StatusUnauthorized != r.StatusCode {
			t.Errorf("expected status code: %v; got %v", http.StatusUnauthorized, r.StatusCode)
		}

		t.Logf("Test_Router tabletest: %d, name: %s", testCount, test.name)
		testCount++
	}
}

func Test_Post_Todos_With_Id(t *testing.T) {
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	var todo []Todo
	db.Preload("Tasks").First(&todo)
	//create POST request based on exiting todo with ids
	json_content, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("create json error: %v", err)
	}
	body := bytes.NewReader(json_content)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/todos?key=%s", srv.URL, TODOLIST_API_KEY), body)
	if err != nil {
		t.Fatalf("create request error: %v", err)
	}
	// send request
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("send request error: %v", err)
	}
	defer r.Body.Close()

	// check status code
	if http.StatusBadRequest != r.StatusCode {
		t.Errorf("expected status code: %v; got %v", http.StatusBadRequest, r.StatusCode)
	}
	t.Logf("Test_Post_Todos_With_Id test: %d, name: post todo with existing ids", testCount)
	testCount++
}

func Test_Post_Todos(t *testing.T) {
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	//create request based on json_todo_post
	json_todo := strings.NewReader(json_todo_post)
	requ, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/todos?key=%s", srv.URL, TODOLIST_API_KEY), json_todo)
	if err != nil {
		t.Fatalf("create request error: %v", err)
	}
	// send request
	resp, err := client.Do(requ)
	if err != nil {
		t.Fatalf("send request error: %v", err)
	}
	defer resp.Body.Close()

	// check status code
	if http.StatusOK != resp.StatusCode {
		t.Errorf("expected status code: %v; got %v", http.StatusOK, resp.StatusCode)
	}
	var response_todo Todo
	json.NewDecoder(resp.Body).Decode(&response_todo)
	if 0 >= response_todo.ID {
		t.Error("expected json with id")
	}
	// get expected content
	var updated_todo Todo
	db.Preload("Tasks").First(&updated_todo, response_todo.ID)
	// normalize content for comparison
	json_expect, err := json.Marshal(updated_todo)
	if err != nil {
		t.Fatalf("create expected json error: %v", err)
	}
	json_response, err := json.Marshal(response_todo)
	if err != nil {
		t.Fatalf("create response json error: %v", err)
	}
	// compare content as string
	expected := string(json_expect)
	response := string(json_response)
	if response != expected {
		t.Errorf("expected content: %v;\ngot %v\n", expected, response)
	}
	t.Logf("Test_Post_Todos test: %d, name: post todo", testCount)
	testCount++
}

func Test_Get_Todo(t *testing.T) {
	//init postgres db
	//init_db()
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	var todos []Todo
	db.Preload("Tasks").Find(&todos)
	for _, todo := range todos {
		//create GET request based on todo.id
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/todos/%d?key=%s", srv.URL, todo.ID, TODOLIST_API_KEY), nil)
		if err != nil {
			t.Fatalf("create request error: %v", err)
		}

		// send request
		r, err := client.Do(req)
		if err != nil {
			t.Fatalf("send request error: %v", err)
		}
		defer r.Body.Close()

		// check status code
		if http.StatusOK != r.StatusCode {
			t.Errorf("expected status code: %v; got %v", http.StatusOK, r.StatusCode)
		}
		// get content
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("get response body error: %v", err)
		}
		// get expected content
		json_expect, err := json.Marshal(todo)
		if err != nil {
			t.Fatalf("create expected json error: %v", err)
		}
		// normalize content for comparison
		compact_repsonse := new(bytes.Buffer)
		compact_expected := new(bytes.Buffer)
		err = json.Compact(compact_repsonse, content)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Compact(compact_expected, json_expect)
		if err != nil {
			fmt.Println(err)
		}
		// compare content as string
		response := compact_repsonse.String()
		expected := compact_expected.String()
		if response != expected {
			t.Errorf("expected content: %v;\ngot %v\n", expected, response)
		}

		t.Logf("Test_Get_Todo: %d, name: GET %s/todos/%d", testCount, srv.URL, todo.ID)
		testCount++
	}
}

func Test_Put_Todo(t *testing.T) {
	//init postgres db
	//init_db()
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	var todos []Todo
	db.Preload("Tasks").Find(&todos)
	for _, todo := range todos {
		//modify todo
		todo.Name = "Update todo name with PUT"
		todo.Tasks = append(todo.Tasks, Task{Name: "Update task name with PUT"})
		//create PUT request based on modified todo
		json_content, err := json.Marshal(todo)
		if err != nil {
			t.Fatalf("create json error: %v", err)
		}
		body := bytes.NewReader(json_content)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/todos/%d?key=%s", srv.URL, todo.ID, TODOLIST_API_KEY), body)
		if err != nil {
			t.Fatalf("create request error: %v", err)
		}

		// send request
		r, err := client.Do(req)
		if err != nil {
			t.Fatalf("send request error: %v", err)
		}
		defer r.Body.Close()

		// check status code
		if http.StatusOK != r.StatusCode {
			t.Errorf("expected status code: %v; got %v", http.StatusOK, r.StatusCode)
		}
		// get content
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("get response body error: %v", err)
		}
		// get expected content
		var updated_todo Todo
		db.Preload("Tasks").First(&updated_todo, todo.ID)
		json_expect, err := json.Marshal(updated_todo)
		if err != nil {
			t.Fatalf("create expected json error: %v", err)
		}
		// normalize content for comparison
		compact_repsonse := new(bytes.Buffer)
		compact_expected := new(bytes.Buffer)
		err = json.Compact(compact_repsonse, content)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Compact(compact_expected, json_expect)
		if err != nil {
			fmt.Println(err)
		}
		// compare content as string
		response := compact_repsonse.String()
		expected := compact_expected.String()
		if response != expected {
			t.Errorf("expected content: %v;\ngot %v\n", expected, response)
		}

		t.Logf("Test_Put_Todo: %d, name: PUT %s/todos/%d", testCount, srv.URL, todo.ID)
		testCount++
	}
}

func Test_Post_Todo(t *testing.T) {
	//init postgres db
	//init_db()
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	var todos []Todo
	db.Preload("Tasks").Find(&todos)
	for _, todo := range todos {
		//modify todo
		todo.Name = "Update todo name with POST"
		todo.Tasks = append(todo.Tasks, Task{Name: "Update task name with POST"})
		//create PUT request based on modified todo
		json_content, err := json.Marshal(todo)
		if err != nil {
			t.Fatalf("create json error: %v", err)
		}
		body := bytes.NewReader(json_content)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/todos/%d?key=%s", srv.URL, todo.ID, TODOLIST_API_KEY), body)
		if err != nil {
			t.Fatalf("create request error: %v", err)
		}

		// send request
		r, err := client.Do(req)
		if err != nil {
			t.Fatalf("send request error: %v", err)
		}
		defer r.Body.Close()

		// check status code
		if http.StatusMethodNotAllowed != r.StatusCode {
			t.Errorf("expected status code: %v; got %v", http.StatusMethodNotAllowed, r.StatusCode)
		}

		t.Logf("Test_Post_Todo not allowed: %d, name: POST %s/todos/%d", testCount, srv.URL, todo.ID)
		testCount++
	}
}

func Test_Patch_Todo(t *testing.T) {
	//init postgres db
	//init_db()
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	var todos []Todo
	db.Preload("Tasks").Find(&todos)
	for _, todo := range todos {
		//modify todo
		todo.Name = "Update todo name with PATCH"
		todo.Tasks = append(todo.Tasks, Task{Name: "Update task name with PATCH"})
		//create PUT request based on modified todo
		json_content, err := json.Marshal(todo)
		if err != nil {
			t.Fatalf("create json error: %v", err)
		}
		body := bytes.NewReader(json_content)
		req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/todos/%d?key=%s", srv.URL, todo.ID, TODOLIST_API_KEY), body)
		if err != nil {
			t.Fatalf("create request error: %v", err)
		}

		// send request
		r, err := client.Do(req)
		if err != nil {
			t.Fatalf("send request error: %v", err)
		}
		defer r.Body.Close()

		// check status code
		if http.StatusMethodNotAllowed != r.StatusCode {
			t.Errorf("expected status code: %v; got %v", http.StatusMethodNotAllowed, r.StatusCode)
		}

		t.Logf("Test_Patch_Todo not allowed: %d, name: PATCH %s/todos/%d", testCount, srv.URL, todo.ID)
		testCount++
	}
}

func Test_Delete_Todo(t *testing.T) {
	//init postgres db
	//init_db()
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	var todos []Todo
	db.Preload("Tasks").Find(&todos)
	for _, todo := range todos {
		//create Delete request based on todo.id
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/todos/%d?key=%s", srv.URL, todo.ID, TODOLIST_API_KEY), nil)
		if err != nil {
			t.Fatalf("create request error: %v", err)
		}

		// send request
		r, err := client.Do(req)
		if err != nil {
			t.Fatalf("send request error: %v", err)
		}
		defer r.Body.Close()

		// check status code
		if http.StatusOK != r.StatusCode {
			t.Errorf("expected status code: %v; got %v", http.StatusOK, r.StatusCode)
		}
		// get content
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("get response body error: %v", err)
		}
		// get expected content
		json_expect, err := json.Marshal(todo)
		if err != nil {
			t.Fatalf("create expected json error: %v", err)
		}
		// normalize content for comparison
		compact_repsonse := new(bytes.Buffer)
		compact_expected := new(bytes.Buffer)
		err = json.Compact(compact_repsonse, content)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Compact(compact_expected, json_expect)
		if err != nil {
			fmt.Println(err)
		}
		// compare content as string
		response := compact_repsonse.String()
		expected := compact_expected.String()
		if response != expected {
			t.Errorf("expected content: %v;\ngot %v\n", expected, response)
		}

		t.Logf("Test_Delete_Todo: %d, name: DELETE %s/todos/%d", testCount, srv.URL, todo.ID)
		testCount++
	}
}

func Test_Delete_Todos(t *testing.T) {
	//init postgres db
	//init_db()
	srv := httptest.NewServer(muxRouter())
	client := srv.Client()
	defer srv.Close()
	var todos []Todo
	db.Unscoped().Where("deleted_at is NOT NULL").Preload("Tasks").Find(&todos)
	//create Delete request
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/todos?key=%s", srv.URL, TODOLIST_API_KEY), nil)
	if err != nil {
		t.Fatalf("create request error: %v", err)
	}

	// send request
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("send request error: %v", err)
	}
	defer r.Body.Close()

	// check status code
	if http.StatusOK != r.StatusCode {
		t.Errorf("expected status code: %v; got %v", http.StatusOK, r.StatusCode)
	}
	// get content
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("get response body error: %v", err)
	}
	// get expected content
	json_expect, err := json.Marshal(todos)
	if err != nil {
		t.Fatalf("create expected json error: %v", err)
	}
	// normalize content for comparison
	compact_repsonse := new(bytes.Buffer)
	compact_expected := new(bytes.Buffer)
	err = json.Compact(compact_repsonse, content)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Compact(compact_expected, json_expect)
	if err != nil {
		fmt.Println(err)
	}
	// compare content as string
	response := compact_repsonse.String()
	expected := compact_expected.String()
	if response != expected {
		t.Errorf("expected content: %v;\ngot %v\n", expected, response)
	}
	t.Logf("Test_Delete_Todo: %d, name: DELETE %s/todos", testCount, srv.URL)
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
