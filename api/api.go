package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/heyjoakim/devops-21/api/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/heyjoakim/devops-21/models"
	"golang.org/x/crypto/bcrypt"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
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

// QueryDb queries the database and returns a list of dictionaries.
func queryDb(query string, args ...interface{}) *sql.Rows {
	liteDB, _ := sql.Open("sqlite3", database)

	res, err := liteDB.Query(query, args...)
	if err != nil {
		log.Print("Database error:", err)
	}
	liteDB.Close()
	return res
}

func serialize(input interface{}) ([]byte, error) {
	js, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func getQueryParameter(r *http.Request, name string) string {
	queryValue := r.URL.Query().Get(name)
	return queryValue
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

// getUserID returns user ID for username
func getUserID(username string) (int, error) {
	var ID int
	err := db.QueryRow("select user_id from user where username = ?", username).Scan(&ID)
	return ID, err
}
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

// Messages godoc
// @Summary Registers a user
// @Description Registers a user, provided that the given info passes all checks.
// @Accept json
// @Produce json
// @Success 203
// @Failure 400 {string} string "unauthorized"
// @Router /register [post]
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var registerRequest models.RegisterRequest
	json.NewDecoder(r.Body).Decode(&registerRequest)

	updateLatest(r)
	var registerError string
	if r.Method == "POST" {
		if len(registerRequest.Username) == 0 {
			registerError = "You have to enter a username"
		} else if len(registerRequest.Email) == 0 || strings.Contains(registerRequest.Email, "@") == false {
			registerError = "You have to enter a valid email address"
		} else if len(registerRequest.Password) == 0 {
			registerError = "You have to enter a password"
		} else if _, err := getUserID(registerRequest.Username); err == nil {
			registerError = "The username is already taken"
		} else {
			statement, err := db.Prepare(`insert into user (username, email, pw_hash) values(?,?,?)`)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Print(err)
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Print(err)
			}

			_, err = statement.Exec(registerRequest.Username, registerRequest.Email, hash)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			statement.Close()

		}
		if registerError != "" {
			error := map[string]interface{}{ // could also be an array
				"status":    400,
				"error_msg": registerError,
			}
			json, _ := serialize(error)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(json)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

	}
}

// Messages godoc
// @Summary Gets the latest messages
// @Description Gets the latest messages in descending order.
// @Param no query int false "Number of results returned"
// @Produce  json
// @Success 200 {object} interface{}
// @Failure 401 {string} string "unauthorized"
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
	noMsgsQuery := getQueryParameter(r, "no")
	if noMsgsQuery == "" {
		noMsgs = 100
	} else {
		noMsgs, _ = strconv.Atoi(noMsgsQuery)
	}

	if r.Method == "GET" {
		query := "SELECT message.*, user.* FROM message, user " +
			"WHERE message.flagged = 0 AND message.author_id = user.user_id " +
			"ORDER BY message.pub_date DESC LIMIT ?"

		messages := queryDb(query, noMsgs)
		var filteredMsgs []models.MessageResponse
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
			err := messages.Scan(&messageID, &authorID, &text, &pubDate, &flagged, &userID, &username, &email, &pwHash)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			filteredMsgs = append(filteredMsgs, models.MessageResponse{
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

// MessagesPerUser godoc
// @Summary Gets the latest messages per user
// @Description Gets the latest messages per user
// @Param no query int false "Number of results returned"
// @Param latest query int false "Something about latest"
// @Produce  json
// @Success 200 {object} interface{}
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string response.Error
// @Router /msgs/{username} [get]
// @Router /msgs/{username} [post]
func messagesPerUserHandler(w http.ResponseWriter, r *http.Request) {
	updateLatest(r)
	params := mux.Vars(r)
	username := params["username"]

	userID, err := getUserID(username)

	notFromSimResponse := notReqFromSimulator(r)
	if notFromSimResponse != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(notFromSimResponse)
		return
	}

	var noMsgs int
	noMsgsQuery := getQueryParameter(r, "no")
	if noMsgsQuery == "" {
		noMsgs = 100
	} else {
		noMsgs, _ = strconv.Atoi(noMsgsQuery)
	}

	if r.Method == "GET" {
		query := "SELECT message.*, user.* FROM message, user " +
			"WHERE message.flagged = 0 AND " +
			"user.user_id = message.author_id AND user.user_id = ?" +
			"ORDER BY message.pub_date DESC LIMIT ?"

		messages := queryDb(query, userID, noMsgs)
		filteredMsgs := make([]models.MessageResponse, 0)
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

			filteredMsgs = append(filteredMsgs, models.MessageResponse{
				Content: text,
				PubDate: pubDate,
				User:    username,
			})
		}
		json, _ := serialize(filteredMsgs)

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)

	} else if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var jsonBody map[string]string
		err := decoder.Decode(&jsonBody)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		statement, _ := db.Prepare("INSERT INTO message (author_id, text, pub_date, flagged) " +
			"VALUES (?, ?, ?, 0)")
		if text := jsonBody["content"]; text == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = statement.Exec(userID, jsonBody["content"], time.Now().Unix())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}

}

// Follow godoc
// @Summary Follow, unfollow or get followers
// @Description Eiter follows a user, unfollows a user or returns a list of users's followers
// @Param no query int false "Number of results returned"
// @Param latest query int false "Something about latest"
// @Accept  json
// @Produce json
// @Success 200 {object} interface{}
// @Success 204 {object} interface{}
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string response.Error
// @Router /fllws/{username} [get]
// @Router /fllws/{username} [post]
func followHandler(w http.ResponseWriter, r *http.Request) {
	updateLatest(r)

	notFromSimResponse := notReqFromSimulator(r)
	if notFromSimResponse != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(notFromSimResponse)
		return
	}

	username := mux.Vars(r)["username"]
	userID, err := getUserID(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %s", username), http.StatusNotFound)
		return
	}

	var followRequest models.FollowRequest
	json.NewDecoder(r.Body).Decode(&followRequest)
	if followRequest.Follow != "" && followRequest.Unfollow != "" {
		http.Error(w, "Invalid input. Can ONLY handle either follow OR unfollow.", http.StatusUnprocessableEntity)
	} else if r.Method == "POST" && followRequest.Follow != "" {
		followsUserID, err := getUserID(followRequest.Follow)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			log.Fatal(err)
		}

		stmt, err := db.Prepare("INSERT INTO follower (who_id, whom_id) VALUES (?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal()
		}

		_, err = stmt.Exec(userID, followsUserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusNoContent)

	} else if r.Method == "POST" && followRequest.Unfollow != "" {
		unfollowsUserID, err := getUserID(followRequest.Unfollow)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			log.Fatal(err)
		}

		stmt, err := db.Prepare("DELETE FROM follower WHERE who_id=? and WHOM_ID=?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal()
		}

		_, err = stmt.Exec(userID, unfollowsUserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusNoContent)

	} else if r.Method == "GET" {
		var noFollowers int
		noFollowersQuery := getQueryParameter(r, "no")
		if noFollowersQuery == "" {
			noFollowers = 100
		} else {
			noFollowers, _ = strconv.Atoi(noFollowersQuery)
		}

		res, _ := db.Query("SELECT user.username FROM user "+
			"INNER JOIN follower ON follower.whom_id=user.user_id "+
			"WHERE follower.who_id=? "+
			"LIMIT ?", userID, noFollowers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		followers := make([]string, 0)
		for res.Next() {
			var username string
			err := res.Scan(&username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			followers = append(followers, username)
		}

		json, _ := serialize(map[string]interface{}{"follows": followers})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
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

	// Swagger
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8001", router))

}
