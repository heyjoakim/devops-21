package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/metrics"
	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
)

// MessagesHandler godoc
// @Summary Gets the latest messages
// @Description Gets the latest messages in descending order.
// @Param no query int false "Number of results returned"
// @Produce  json
// @Success 200 {object} interface{}
// @Failure 401 {string} string "unauthorized"
// @Router /api/msgs [get]
func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	hist := metrics.GetHistogramVec("get_api_msgs")
	if hist != nil {
		timer := createEndpointTimer(hist)
		defer timer.ObserveDuration()
	}

	updateLatest(r)

	notFromSimResponse := helpers.IsFromSimulator(w, r)
	if notFromSimResponse {
		return
	}

	var noMsgs int
	noMsgsQuery := helpers.GetQueryParameter(r, "no")
	if noMsgsQuery == "" {
		noMsgs = 100
	} else {
		noMsgs, _ = strconv.Atoi(noMsgsQuery)
	}

	results := services.GetPublicMessages(noMsgs)

	jsonData, _ := helpers.Serialize(results)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonData)
}

// GetMessagesFromUserHandler godoc
// @Summary Gets the latest messages from a specific user
// @Description Gets the latest messages in descending order from a specific user.
// @Param username query string true "Username"
// @Param no query int false "Number of results returned"
// @Produce  json
// @Success 200 {object} interface{}
// @Failure 401 {string} string "unauthorized"
// @Router /api/msgs/{username} [get]
func GetMessagesFromUserHandler(w http.ResponseWriter, r *http.Request) {
	hist := metrics.GetHistogramVec("get_api_msgs_username")
	if hist != nil {
		timer := createEndpointTimer(hist)
		defer timer.ObserveDuration()
	}

	updateLatest(r)

	params := mux.Vars(r)
	username := params["username"]

	userID, _ := services.GetUserID(username)

	notFromSimResponse := helpers.IsFromSimulator(w, r)
	if notFromSimResponse {
		return
	}

	var noMsgs int
	noMsgsQuery := helpers.GetQueryParameter(r, "no")
	if noMsgsQuery == "" {
		noMsgs = 100
	} else {
		noMsgs, _ = strconv.Atoi(noMsgsQuery)
	}

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
	_, _ = w.Write(jsonData)
}

// PostMessageHandler godoc
// @Summary Create a message from user
// @Description Creates a message from a specific user.
// @Param username query string true "Username"
// @Produce json
// @Success 200 {object} interface{}
// @Failure 401 {string} string "unauthorized"
// @Router /api/msgs/{username} [post]
func PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	hist := metrics.GetHistogramVec("post_api_msgs_username")
	if hist != nil {
		timer := createEndpointTimer(hist)
		defer timer.ObserveDuration()
	}

	updateLatest(r)
	params := mux.Vars(r)
	username := params["username"]

	userID, _ := services.GetUserID(username)

	notFromSimResponse := helpers.IsFromSimulator(w, r)
	if notFromSimResponse {
		return
	}

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
