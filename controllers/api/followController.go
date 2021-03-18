package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	"log"
	"net/http"
	"strconv"
)

// FollowHandler godoc
// @Summary Follows a user, unfollows a user
// @Description follows a user, unfollows a user
// @Param latest query int false "Something about latest"
// @Accept  json
// @Success 204 {object} interface{}
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string response.Error
// @Router /api/fllws/{username} [post]
func FollowHandler(w http.ResponseWriter, r *http.Request) {
	updateLatest(r)

	notFromSimResponse := helpers.NotReqFromSimulator(r)
	if notFromSimResponse != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write(notFromSimResponse)
		return
	}

	username := mux.Vars(r)["username"]
	userID, err := services.GetUserID(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %s", username), http.StatusNotFound)
		return
	}

	var followRequest models.FollowRequest
	_ = json.NewDecoder(r.Body).Decode(&followRequest)
	if followRequest.Follow == "" || followRequest.Unfollow == "" {
		if followRequest.Follow != "" {
			followsUserID, err := services.GetUserID(followRequest.Follow)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			follower := models.Follower{WhoID: userID, WhomID: followsUserID}
			err = services.CreateFollower(follower)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusNoContent)
		} else if followRequest.Unfollow != "" {
			unfollowsUserID, err := services.GetUserID(followRequest.Unfollow)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				log.Fatal(err)
			}

			err = services.UnfollowUser(userID, unfollowsUserID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal()
			}

			w.WriteHeader(http.StatusNoContent)
		}
	} else {
		http.Error(w, "Invalid input. Can ONLY handle either follow OR unfollow.", http.StatusUnprocessableEntity)
	}
}

// GetFollowersHandler godoc
// @Summary Get followers
// @Description Returns a list of users followers
// @Param no query int false "Number of results returned"
// @Param latest query int false "Something about latest"
// @Accept  json
// @Produce json
// @Success 200 {object} interface{}
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string response.Error
// @Router /api/fllws/{username} [get]
func GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
	updateLatest(r)

	notFromSimResponse := helpers.NotReqFromSimulator(r)
	if notFromSimResponse != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write(notFromSimResponse)
		return
	}

	username := mux.Vars(r)["username"]
	userID, err := services.GetUserID(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %s", username), http.StatusNotFound)
		return
	}

	var noFollowers int
	noFollowersQuery := helpers.GetQueryParameter(r, "no")
	if noFollowersQuery == "" {
		noFollowers = 100
	} else {
		noFollowers, _ = strconv.Atoi(noFollowersQuery)
	}

	users := services.GetAllUsersFollowers(userID, noFollowers)

	jsonData, _ := helpers.Serialize(map[string]interface{}{"follows": users})
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonData)
}
