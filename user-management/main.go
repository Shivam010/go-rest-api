// A Simple Demo GoLang HTTP JSON API for User-Management
// Capable of Creating a user, Returning it, Editing it, and Deleting it.

// Database Structure
// DBMS: "PostgreSQL"
// Schema: "user_management"
// Table: "users"
// Columns: "id serial, fname text, lname text, dob text, email text, phono_no bigint"

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

const (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "postgres"
	dbPassword = "appointy"
	dbName     = "test"
)

// User Object
type User struct {
	ID      int64  `json:"id"`
	Fname   string `json:"fname"`
	Lname   string `json:"lname"`
	DOB     string `json:"dob"`
	Email   string `json:"email"`
	PhoneNo int64  `json:"phoneno"`
}

type empty struct {
}

// RequestHandlerFunc is the type defined to use the http Handler Function externally
type RequestHandlerFunc func(http.ResponseWriter, *http.Request)

// wrapper wraps the http request with sequence of middlewares provided
func wrapper(fn RequestHandlerFunc, mds ...func(RequestHandlerFunc) RequestHandlerFunc) RequestHandlerFunc {
	for _, md := range mds {
		fn = md(fn)
	}
	return fn
}

// BasicAuthentication middleware
func BasicAuthentication(req RequestHandlerFunc) RequestHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic Realm: "Restricted"`)
		user, pass, ok := r.BasicAuth()
		if !ok || (ok && (user != "mavis" || pass != "shivam")) {
			http.Error(w, "Unauthorized Access", http.StatusUnauthorized)
			return
		}
		req(w, r)
	}
}

// CreateUser create user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error decoding data: %v", err)
		return
	}
	fmt.Println(user)
	const query = `INSERT INTO user_management.users(fname, lname, dob, email, phone_no) VALUES($1, $2, $3, $4, $5) returning id`
	err = db.QueryRow(query, user.Fname, user.Lname, user.DOB, user.Email, user.PhoneNo).Scan(&user.ID)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("create user error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error encoding data: %v", err)
		return
	}
}

// GetUser returns a user
func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error in string conversion: %v", err)
		return
	}
	const query = `SELECT * FROM user_management.users WHERE id = $1`
	row := db.QueryRow(query, id)
	user := User{}
	row.Scan(&user.Fname, &user.Lname, &user.DOB, &user.Email, &user.PhoneNo, &user.ID)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error encoding data: %v", err)
		return
	}
}

// GetAllUser returns all user
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	const query = `SELECT * FROM user_management.users`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error: %v", err)
		return
	}
	list := []User{}
	for rows.Next() {
		user := User{}
		rows.Scan(&user.Fname, &user.Lname, &user.DOB, &user.Email, &user.PhoneNo, &user.ID)
		list = append(list, user)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error encoding data: %v", err)
		return
	}
}

// EditUser edit a user
func EditUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error in string conversion: %v", err)
		return
	}
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error encoding data: %v", err)
		return
	}
	user.ID = id
	fmt.Println(user)
	const query = `UPDATE user_management.users SET fname=$2, lname=$3, dob=$4, email=$5, phone_no=$6 WHERE id = $1;`
	_, err = db.Exec(query, user.ID, user.Fname, user.Lname, user.DOB, user.Email, user.PhoneNo)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(empty{}); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error encoding data: %v", err)
		return
	}
}

// DeleteUser deletes a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error in string conversion: %v", err)
		return
	}
	const query = `DELETE FROM user_management.users WHERE id = $1`
	if _, err = db.Exec(query, id); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(empty{}); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error encoding data: %v", err)
		return
	}
}

func main() {
	// database connection
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
		return
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping error: %v", err)
		return
	}
	defer db.Close()

	// api pattern handlers
	http.HandleFunc("/create", wrapper(CreateUser, BasicAuthentication)) // POST
	http.HandleFunc("/user", wrapper(GetUser, BasicAuthentication))      // GET
	http.HandleFunc("/users", wrapper(GetAllUser, BasicAuthentication))  // GET
	http.HandleFunc("/edit", wrapper(EditUser, BasicAuthentication))     // PUT
	http.HandleFunc("/delete", wrapper(DeleteUser, BasicAuthentication)) // DELETE

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
