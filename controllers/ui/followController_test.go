package ui

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	"github.com/stretchr/testify/assert"
)

func TestMemoryFollow(t *testing.T) {
	// Setup two mock users
	foo := &models.User{Username: "Foo", Email: "foo@baz.com", PwHash: "off"}
	bar := &models.User{Username: "Bar", Email: "bar@baz.com", PwHash: "rab"}
	fooData := url.Values{"username": {foo.Username},
		"email":     {foo.Email},
		"password":  {foo.PwHash},
		"password2": {foo.PwHash}}
	barData := url.Values{"username": {bar.Username},
		"email":     {bar.Email},
		"password":  {bar.PwHash},
		"password2": {bar.PwHash}}

	MemoryTimelineHelper(
		fooData,
		url.Values{
			"text":  {"Foo test message"},
			"token": {foo.Username},
		},
		barData,
		url.Values{
			"text":  {"Bar test message"},
			"token": {bar.Username},
		},
	)

	// Sequp request to follow user foo
	req, _ := http.NewRequest("POST", "/username/follow", nil)
	w := httptest.NewRecorder()

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set current user to bar
	session, _ := store.Get(req, "_cookie")
	session.Values["user_id"], _ = services.GetUserID(bar.Username)
	session.Values["username"] = bar.Username
	_ = session.Save(req, w)

	// Set URL vars to be retrieved by mux.Vars
	// Expected params["username"] = foo
	req = mux.SetURLVars(req, map[string]string{"username": foo.Username})

	handler := http.HandlerFunc(FollowUserHandler)
	handler.ServeHTTP(w, req)
	resp := w.Result()
	assert.Equal(t, resp.StatusCode, 302)

	// Beigin new request to check the users page for the added message
	newReq, _ := http.NewRequest("GET", "/username", nil)
	newW := httptest.NewRecorder()
	newReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set URL vars to be retrieved by mux.Vars, session already set to bar
	newReq = mux.SetURLVars(newReq, map[string]string{"username": bar.Username})

	// New request, need to set session vars again
	session, _ = store.Get(newReq, "_cookie")
	session.Values["user_id"], _ = services.GetUserID(bar.Username)
	session.Values["username"] = bar.Username
	_ = session.Save(newReq, newW)

	// Handle request
	newHandler := http.HandlerFunc(UserTimelineHandler)
	newHandler.ServeHTTP(newW, newReq)
	newResp := newW.Result()
	newBody, _ := ioutil.ReadAll(newResp.Body)

	// Assert that new message is added to the personal page
	assert.Contains(t, string(newBody), "Foo test message")
	assert.Contains(t, string(newBody), "Bar test message")
}
