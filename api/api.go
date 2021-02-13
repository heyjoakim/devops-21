package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	database = "../tmp/minitwit.db"
	db       *sql.DB
	latest   = 0
	debug    = true
)

// connectDb returns a new connection to the database.
func connectDb() (*sql.DB, error) {
	return sql.Open("sqlite3", database)
}

func serialize(input interface{}) ([]byte, error) {
	js, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func notReqFromSimulator(r *http.Request) {}

func getUserID(username string) string { return "" }

func beforeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Insert here
		next.ServeHTTP(w, r)
	})
}

func afterRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		// Insert here
	})
}
func updateLatest(r *http.Request) {}

func getLatestHandler(w http.ResponseWriter, r *http.Request) {}

func registerHandler(w http.ResponseWriter, r *http.Request) {}

func messagesHandler(w http.ResponseWriter, r *http.Request) {}

func messagesPerUserHandler(w http.ResponseWriter, r *http.Request) {}

func followHandler(w http.ResponseWriter, r *http.Request) {}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{ // could also be an array
		"ping": "pong",
	}

	// could also be an array
	// data := [1]string{"pong"}

	json, err := serialize(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // if WriteHeader is not set, http status code will be 200 OK
	w.Write(json)
}

// init is automatically executed on program startup. Can't be called
// or referenced.
func init() {
	database, err := connectDb()
	if err != nil {
		log.Fatal(err)
	}
	db = database
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/latest", getLatestHandler)
	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/msgs", messagesHandler)
	router.HandleFunc("/msgs/{username}", messagesPerUserHandler).Methods("GET", "POST")
	router.HandleFunc("/fllws/{username}", followHandler).Methods("GET", "POST")
	router.HandleFunc("/example", exampleHandler)

	log.Fatal(http.ListenAndServe(":8001", router))
}