package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func updateLatest(r *http.Request) {
	tryLatestQuery := r.URL.Query().Get("latest")

	if tryLatestQuery == "" {
		latest = -1
	} else {
		tryLatest, _ := strconv.Atoi(tryLatestQuery)
		latest = tryLatest
	}
}

var (
	latest = 0
)

// GetLatest godoc
// @Summary Get the latest x
// @Description Get the latest x
// @Produce  json
// @Success 200 {object} interface{}
// @Router /api/latest [get]
func GetLatestHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{ // could also be an array
		"latest": latest,
	}

	jsonData, err := helpers.Serialize(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// Messages godoc
// @Summary Registers a user
// @Description Registers a user, provided that the given info passes all checks.
// @Accept json
// @Produce json
// @Success 203
// @Failure 400 {string} string "unauthorized"
// @Router /api/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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
		} else if _, err := services.GetUserID(registerRequest.Username); err == nil {
			registerError = "The username is already taken"
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Print(err)
			}

			user := models.User{Username: registerRequest.Username, Email: registerRequest.Email, PwHash: string(hash)}
			err = services.CreateUser(user)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if registerError != "" {
			error := map[string]interface{}{
				"status":    400,
				"error_msg": registerError,
			}
			jsonData, _ := helpers.Serialize(error)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonData)
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
// @Router /api/msgs [get]
func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	updateLatest(r)

	notFromSimResponse := helpers.NotReqFromSimulator(r)
	if notFromSimResponse != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(notFromSimResponse)
		return
	}

	var noMsgs int
	noMsgsQuery := helpers.GetQueryParameter(r, "no")
	if noMsgsQuery == "" {
		noMsgs = 100
	} else {
		noMsgs, _ = strconv.Atoi(noMsgsQuery)
	}

	if r.Method == "GET" {
		results := services.GetPublicMessages(noMsgs)

		jsonData, _ := helpers.Serialize(results)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
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
// @Router /api/msgs/{username} [get]
// @Router /api/msgs/{username} [post]
func MessagesPerUserHandler(w http.ResponseWriter, r *http.Request) {
	updateLatest(r)
	params := mux.Vars(r)
	username := params["username"]

	userID, _ := services.GetUserID(username)

	notFromSimResponse := helpers.NotReqFromSimulator(r)
	if notFromSimResponse != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(notFromSimResponse)
		return
	}

	var noMsgs int
	noMsgsQuery := helpers.GetQueryParameter(r, "no")
	if noMsgsQuery == "" {
		noMsgs = 100
	} else {
		noMsgs, _ = strconv.Atoi(noMsgsQuery)
	}

	if r.Method == "GET" {
		results := services.GetMessagesForUser(noMsgs, userID)

		var messages []models.MessageResponse
		for _, result := range results {
			message := models.MessageResponse{
				User:    result.Username,
				Content: result.Content,
				PubDate: result.PubDate,
			}
			messages = append(messages, message)
		}

		jsonData, _ := helpers.Serialize(messages)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)

	} else if r.Method == "POST" {
		var messageRequest models.MessageRequest
		err := json.NewDecoder(r.Body).Decode(&messageRequest)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		message := models.Message{AuthorID: userID, Text: messageRequest.Content, PubDate: time.Now().Unix(), Flagged: 0}
		err = services.CreateMessage(message)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}

}

//// FollowHandler godoc
//// @Summary Follow, unfollow or get followers
//// @Description Eiter follows a user, unfollows a user or returns a list of users's followers
//// @Param no query int false "Number of results returned"
//// @Param latest query int false "Something about latest"
//// @Accept  json
//// @Produce json
//// @Success 200 {object} interface{}
//// @Success 204 {object} interface{}
//// @Failure 401 {string} string "unauthorized"
//// @Failure 500 {string} string response.Error
//// @Router /api/fllws/{username} [get]
//// @Router /api/fllws/{username} [post]
//func FollowHandler(w http.ResponseWriter, r *http.Request) {
//	updateLatest(r)
//
//	notFromSimResponse := helpers.NotReqFromSimulator(r)
//	if notFromSimResponse != nil {
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusUnauthorized)
//		w.Write(notFromSimResponse)
//		return
//	}
//
//	username := mux.Vars(r)["username"]
//	userID, err := services.GetUserID(username)
//	if err != nil {
//		http.Error(w, fmt.Sprintf("User not found: %s", username), http.StatusNotFound)
//		return
//	}
//
//	var followRequest models.FollowRequest
//	json.NewDecoder(r.Body).Decode(&followRequest)
//	if followRequest.Follow != "" && followRequest.Unfollow != "" {
//		http.Error(w, "Invalid input. Can ONLY handle either follow OR unfollow.", http.StatusUnprocessableEntity)
//	} else if r.Method == "POST" && followRequest.Follow != "" {
//		followsUserID, err := services.GetUserID(followRequest.Follow)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusNotFound)
//			return
//		}
//
//		follower := models.Follower{WhoID: userID, WhomID: followsUserID}
//		err = services.CreateFollower(follower)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//		}
//		w.WriteHeader(http.StatusNoContent)
//		return
//	} else if r.Method == "POST" && followRequest.Unfollow != "" {
//		unfollowsUserID, err := services.GetUserID(followRequest.Unfollow)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusNotFound)
//			log.Fatal(err)
//		}
//
//		err = services.UnfollowUser(userID, unfollowsUserID)
//
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			log.Fatal()
//		}
//
//		w.WriteHeader(http.StatusNoContent)
//		return
//
//	} else if r.Method == "GET" {
//		var noFollowers int
//		noFollowersQuery := helpers.GetQueryParameter(r, "no")
//		if noFollowersQuery == "" {
//			noFollowers = 100
//		} else {
//			noFollowers, _ = strconv.Atoi(noFollowersQuery)
//		}
//
//		users := services.GetAllUsersFollowers(userID, noFollowers)
//
//		jsonData, _ := helpers.Serialize(map[string]interface{}{"follows": users})
//		w.Header().Set("Content-Type", "application/json")
//		w.Write(jsonData)
//	}
//}
