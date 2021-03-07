package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	"net/http"
	"strconv"
	"time"
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

	results := services.GetPublicMessages(noMsgs)

	jsonData, _ := helpers.Serialize(results)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetMessagesFromUserHandler(w http.ResponseWriter, r *http.Request) {
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
}

func PostMessageHandler(w http.ResponseWriter, r *http.Request) {
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
