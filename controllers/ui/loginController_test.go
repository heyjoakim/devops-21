package ui

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

var loginUserData url.Values = url.Values{
	"username":  {"Bob"},
	"email":     {"bob@go.com"},
	"password":  {"secret"},
	"password2": {"secret"},
}

var loginUser = &models.User{
	Username: "Bob",
	Email:    "bob@go.com",
	PwHash:   "secret",
}

func TestMemoryLoginHelper(t *testing.T) {
	var resp *http.Response
	var body []byte
	mock := url.Values{}

	// Need to register a user to test error message
	MemoryRegisterHelper(loginUserData)

	// Test wrong username
	mock.Add("username", "wrong_username")
	resp = MemoryLoginHelper(mock)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Unknown username")

	// Test missing password
	mock.Set("username", loginUser.Username)
	mock.Add("password", "wrong_password")
	resp = MemoryLoginHelper(mock)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid password")

	// Test successful login
	mock.Set("password", loginUser.PwHash)
	resp = MemoryLoginHelper(mock)
	assert.Equal(t, resp.StatusCode, 302, "A successful login should redirrect")
}
