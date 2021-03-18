package ui

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

var timelineUserData url.Values = url.Values{
	"username":  {"Ron"},
	"email":     {"ron@go.com"},
	"password":  {"secret"},
	"password2": {"secret"},
}

var timelineUser = &models.User{
	Username: "Ron",
	Email:    "ron@go.com",
	PwHash:   "secret",
}

func TestMemoryTimeline(t *testing.T) {
	// Expected message
	var msg string = "Test message on timeline"
	var msgData url.Values = url.Values{
		"text":  {msg},
		"token": {timelineUser.Username},
	}

	// Add message
	MemoryAddMessageHelper(msgData, timelineUserData)

	// Begin new request to check the public page for message
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set session values
	session, _ := store.Get(req, "_cookie")
	session.Values["user_id"] = timelineUser.UserID
	session.Values["username"] = timelineUser.Username
	_ = session.Save(req, w)

	handler := http.HandlerFunc(PublicTimelineHandler)
	handler.ServeHTTP(w, req)
	checkResp := w.Result()
	body, _ := ioutil.ReadAll(checkResp.Body)

	// Assert that new message is added to the page
	assert.Contains(t, string(body), msg)
}
