// A Simple Demo GoLang HTTP JSON API for ToDo List Management

// Database Structure
// DBMS: "PostgreSQL"
// Schema: "todolist_management"
// Table: "todolist"
// Columns: "id serial, fname text, lname text, dob text, email text, phono_no bigint"

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Shivam010/go-rest-api/todolist-management/lib"

	_ "github.com/lib/pq"
)

// TodoListManagement ...
type TodoListManagement struct {
	c *todolist.Core
}

// NewTodoListManagement ...
func NewTodoListManagement(c *todolist.Core) *TodoListManagement {
	return &TodoListManagement{c}
}

// Ex ...
func (t *TodoListManagement) Ex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if err := t.c.Ex("name"); err != nil {
		InternalServerError(w, err)
		return
	}
	ReturnJSONEncoded(w, "Ex Call")
}

// AddDeleteOrEdit ...
func (t *TodoListManagement) AddDeleteOrEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		t.AddTodoList(w, r)
	} else if r.Method == "DELETE" {
		t.DeleteTodoList(w, r)
	} else if r.Method == "PATCH" {
		t.EditTodoListName(w, r)
	} else {
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
}

// AddTodoList ...
func (t *TodoListManagement) AddTodoList(w http.ResponseWriter, r *http.Request) {
	list := &todolist.TodoList{}
	if err := json.NewDecoder(r.Body).Decode(list); err != nil {
		InternalServerError(w, err)
		return
	}

	obj, err := t.c.AddTodoList(list)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	ReturnJSONEncoded(w, obj)
}

// DeleteTodoList ...
func (t *TodoListManagement) DeleteTodoList(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if err := t.c.DeleteTodoList(id); err != nil {
		InternalServerError(w, err)
		return
	}
	ReturnJSONEncoded(w, empty{})
}

// EditTodoListName ...
func (t *TodoListManagement) EditTodoListName(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	list := &todolist.TodoList{}
	if err := json.NewDecoder(r.Body).Decode(list); err != nil {
		InternalServerError(w, err)
		return
	}
	if err := t.c.EditTodoListName(id, list.Name); err != nil {
		InternalServerError(w, err)
		return
	}
	ReturnJSONEncoded(w, empty{})
}

// AddTodoItem ...
func (t *TodoListManagement) AddTodoItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	type Req struct {
		Lid  int64              `json:"list_id"`
		Item *todolist.TodoItem `json:"item"`
	}
	req := &Req{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		InternalServerError(w, err)
		return
	}
	item, err := t.c.AddTodoItem(req.Lid, req.Item)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	ReturnJSONEncoded(w, item)
}

// DeleteTodoListItem ...
func (t *TodoListManagement) DeleteTodoListItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	if err := t.c.DeleteTodoListItem(id); err != nil {
		InternalServerError(w, err)
		return
	}
	ReturnJSONEncoded(w, empty{})
}

// GetTodoListItem ...
func (t *TodoListManagement) GetTodoListItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	item, err := t.c.GetTodoListItem(id)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	ReturnJSONEncoded(w, item)
}

// UpdateTodoItem ...
func (t *TodoListManagement) UpdateTodoItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	item := &todolist.TodoItem{}
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		InternalServerError(w, err)
		return
	}
	if err := t.c.UpdateTodoItem(item); err != nil {
		InternalServerError(w, err)
		return
	}
	ReturnJSONEncoded(w, empty{})
}

// GetTodoList ...
func (t *TodoListManagement) GetTodoList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	list, err := t.c.GetTodoList(id)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	ReturnJSONEncoded(w, list)
}

func main() {
	// database connection
	db, err := DatabaseConnection()
	if err != nil {
		return
	}
	defer db.Close()

	tdm := NewTodoListManagement(todolist.NewCore(db))

	// api pattern handlers
	http.HandleFunc("/ex", Wrapper(tdm.Ex, BasicAuthentication))
	http.HandleFunc("/todolist", Wrapper(tdm.AddDeleteOrEdit, BasicAuthentication))
	http.HandleFunc("/todolist/addItem", Wrapper(tdm.AddTodoItem, BasicAuthentication))
	http.HandleFunc("/todolist/deleteItem", Wrapper(tdm.DeleteTodoListItem, BasicAuthentication))
	http.HandleFunc("/todolist/getItem", Wrapper(tdm.GetTodoListItem, BasicAuthentication))
	http.HandleFunc("/todolist/updateItem", Wrapper(tdm.UpdateTodoItem, BasicAuthentication))
	http.HandleFunc("/todolist/getList", Wrapper(tdm.GetTodoList, BasicAuthentication))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
