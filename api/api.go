package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/heyjoakim/devops-21/api/docs" // docs is generated by Swag CLI, you have to import it.

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	Content string
	PubDate int
	User    string
}

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

// QueryDb queries the database and returns a list of dictionaries.
func queryDb(query string, args ...interface{}) *sql.Rows {
	liteDB, _ := sql.Open("sqlite3", database)

	res, _ := liteDB.Query(query, args...)

	return res
}

func serialize(input interface{}) ([]byte, error) {
	js, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func getQueryParam(r *http.Request, name string, fallback interface{}) interface{} {
	res := r.URL.Query().Get(name)
	if res == "" {
		return fallback
	}
	return res
}

func notReqFromSimulator(r *http.Request) []byte {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		data := map[string]interface{}{
			"status":    http.StatusForbidden,
			"error_msg": "You are not authorized to use this resource!",
		}

		json, _ := serialize(data)
		return json
	}
	return nil
}

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

func updateLatest(r *http.Request) {
	tryLatestQuery := r.URL.Query().Get("latest")

	if tryLatestQuery == "" {
		latest = -1
	} else {
		tryLatest, _ := strconv.Atoi(tryLatestQuery)
		latest = tryLatest
	}
}

// GetLatest godoc
// @Summary Get the latest x
// @Description Get the latest x
// @Produce  json
// @Success 200 {object} interface{}
// @Router /latest [get]
func getLatestHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{ // could also be an array
		"latest": latest,
	}

	json, err := serialize(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {}

// Messages godoc
// @Summary Gets the latest messages
// @Description Gets the latest messages in descending order.
// @Param no query int false "Number of results returned"
// @Produce  json
// @Success 200 {object} interface{}
// @Failure 403 {string} string "unauthorized"
// @Router /msgs [get]
func messagesHandler(w http.ResponseWriter, r *http.Request) {
	updateLatest(r)

	notFromSimResponse := notReqFromSimulator(r)
	if notFromSimResponse != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(notFromSimResponse)
		return
	}

	var noMsgs int
	var err error
	noMsgsQuery := r.URL.Query().Get("no")
	if noMsgsQuery == "" {
		noMsgs = 100
	} else {
		noMsgs, err = strconv.Atoi(noMsgsQuery)
		if err != nil {
			noMsgs = 100
		}
	}

	if r.Method == "GET" {
		query := "SELECT message.*, user.* FROM message, user " +
			"WHERE message.flagged = 0 AND message.author_id = user.user_id " +
			"ORDER BY message.pub_date DESC LIMIT ?"

		messages := queryDb(query, noMsgs)
		var filteredMsgs []Message
		for messages.Next() {
			var (
				messageID int
				authorID  int
				text      string
				pubDate   int
				flagged   int
				userID    int
				username  string
				email     string
				pwHash    string
			)
			err = messages.Scan(&messageID, &authorID, &text, &pubDate, &flagged, &userID, &username, &email, &pwHash)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			filteredMsgs = append(filteredMsgs, Message{
				Content: text,
				PubDate: pubDate,
				User:    username,
			})
		}
		json, _ := serialize(filteredMsgs)

		w.Header().Set("Content-Type", "application/json")

		w.Write(json)
	}
}

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

// @title Minitwit API
// @version 1.0
// @description This API is comsumed by the simulator for the course DevOps, Software Evolution and Software Maintenance @ ITU Spring 2021
// @termsOfService http://swagger.io/terms/
// @host localhost:8001
// @BasePath /
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/latest", getLatestHandler)
	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/msgs", messagesHandler)
	router.HandleFunc("/msgs/{username}", messagesPerUserHandler).Methods("GET", "POST")
	router.HandleFunc("/fllws/{username}", followHandler).Methods("GET", "POST")
	router.HandleFunc("/example", exampleHandler)

	// Swagger
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8001", router))

}
