package ui

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

var messageUserData url.Values = url.Values{
	"username":  {"Tim"},
	"email":     {"tim@go.com"},
	"password":  {"secret"},
	"password2": {"secret"},
}

var messageUser = &models.User{
	Username: "Tim",
	Email:    "tim@go.com",
	PwHash:   "secret",
}

func TestMemoryAddMessage(t *testing.T) {
	var resp *http.Response

	// Expected message
	var msg string = "Test message personal page"
	var msgData url.Values = url.Values{
		"text":  {msg},
		"token": {messageUser.Username},
	}

	// Add message
	resp = MemoryAddMessageHelper(msgData, messageUserData)

	// Assert that adding a message redirects
	assert.Equal(t, resp.StatusCode, 302)

	// Beigin new request to check the users page for the added message
	req, _ := http.NewRequest("GET", "/username", nil)
	w := httptest.NewRecorder()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set URL vars to be retrieved by mux.Vars
	req = mux.SetURLVars(req, map[string]string{"username": messageUser.Username})

	// Set session values
	session, _ := store.Get(req, "_cookie")
	session.Values["user_id"] = messageUser.UserID
	session.Values["username"] = messageUser.Username
	_ = session.Save(req, w)

	// Handle request
	handler := http.HandlerFunc(UserTimelineHandler)
	handler.ServeHTTP(w, req)
	checkResp := w.Result()
	body, _ := ioutil.ReadAll(checkResp.Body)

	// Assert that new message is added to the personal page
	assert.Contains(t, string(body), msg)
}
