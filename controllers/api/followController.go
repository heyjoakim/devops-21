package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/metrics"
	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	log "github.com/sirupsen/logrus"
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
	hist := metrics.GetHistogramVec("post_api_fllws_username")
	if hist != nil {
		timer := createEndpointTimer(hist)
		defer timer.ObserveDuration()
	}

	updateLatest(r)

	notFromSimResponse := helpers.IsFromSimulator(w, r)
	if notFromSimResponse {
		return
	}

	username := mux.Vars(r)["username"]
	userID, err := services.GetUserID(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %s", username), http.StatusNotFound)
		log.Error(fmt.Sprintf("FollowHandler: User not found: %s", username))
		return
	}

	var followRequest models.FollowRequest
	_ = json.NewDecoder(r.Body).Decode(&followRequest)
	if followRequest.Follow == "" || followRequest.Unfollow == "" {
		if followRequest.Follow != "" {
			follow(followRequest, userID, w)
		} else if followRequest.Unfollow != "" {
			unfollow(followRequest, userID, w)
		}
	} else {
		http.Error(w, "Invalid input. Can ONLY handle either follow OR unfollow.", http.StatusUnprocessableEntity)
	}
}

func follow(followRequest models.FollowRequest, userID uint, w http.ResponseWriter) {
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
}

func unfollow(followRequest models.FollowRequest, userID uint, w http.ResponseWriter) {
	unfollowsUserID, err := services.GetUserID(followRequest.Unfollow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Println(err)
	}

	err = services.UnfollowUser(userID, unfollowsUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}

	w.WriteHeader(http.StatusNoContent)
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
	hist := metrics.GetHistogramVec("get_api_fllws_username")
	if hist != nil {
		timer := createEndpointTimer(hist)
		defer timer.ObserveDuration()
	}

	updateLatest(r)

	notFromSimResponse := helpers.IsFromSimulator(w, r)
	if notFromSimResponse {
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
