// A Simple Demo GoLang HTTP JSON API for User-Management
// Capable of Creating a user, Returning it, Editing it, and Deleting it.

// Database Structure
// DBMS: "PostgreSQL"
// Schema: "public"
// Table: "users"

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
	const query = `INSERT INTO public.users(fname, lname, dob, email, phone_no) VALUES($1, $2, $3, $4, $5) returning id`
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
	id := r.URL.Query().Get("id")
	const query = `SELECT * FROM public.users WHERE id = $1`
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
	const query = `SELECT * FROM public.users`
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
	const query = `UPDATE public.users SET fname=$2, lname=$3, dob=$4, email=$5, phone_no=$6 WHERE id = $1;`
	_, err = db.Exec(query, user.ID, user.Fname, user.Lname, user.DOB, user.Email, user.PhoneNo)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
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
	const query = `DELETE FROM public.users WHERE id = $1`
	if _, err = db.Exec(query, id); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
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
	http.HandleFunc("/create", CreateUser) // POST
	http.HandleFunc("/user", GetUser)      // GET
	http.HandleFunc("/users", GetAllUser)  // GET
	http.HandleFunc("/edit", EditUser)     // PUT
	http.HandleFunc("/delete", DeleteUser) // DELETE

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
