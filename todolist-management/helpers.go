package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-rest-api/todolist-management/lib"
	"log"
	"net/http"
)

// a generic empty struct to return a empty JSON object {} in response
type empty struct {
}

const (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "postgres"
	dbPassword = "appointy"
	dbName     = "test"
)

// InternalServerError is a generic internal server error handler
func InternalServerError(w http.ResponseWriter, err error) {
	if err == todolist.ErrNotFound || err == todolist.ErrItemNotFound {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}
	http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	log.Println(err)
	return
}

// ReturnJSONEncoded is a generic response writer for interfaces in JSON content-type
func ReturnJSONEncoded(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		InternalServerError(w, err)
		return
	}
}

// RequestHandlerFunc is the type defined to use the http Handler Function externally
type RequestHandlerFunc func(http.ResponseWriter, *http.Request)

// Wrapper wraps the http request with sequence of middlewares provided
func Wrapper(fn RequestHandlerFunc, mds ...func(RequestHandlerFunc) RequestHandlerFunc) RequestHandlerFunc {
	for _, md := range mds {
		fn = md(fn)
	}
	return fn
}

// BasicAuthentication middleware of Basic Auth
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

// DatabaseConnection returns a database connection setup
func DatabaseConnection() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping error: %v", err)
		return nil, err
	}
	return db, err
}
