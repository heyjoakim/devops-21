package ui

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

var defaultUserData url.Values = url.Values{ //nolint
	"username":  {"Rob"},
	"email":     {"rob@go.com"},
	"password":  {"secret"},
	"password2": {"secret"},
}

var defaultUser = &models.User{
	Username: "Rob",
	Email:    "rob@go.com",
	PwHash:   "secret",
}

func TestMemoryRegister(t *testing.T) {
	var resp *http.Response
	var body []byte
	mock := url.Values{}

	// // Test missing username
	resp = MemoryRegisterHelper(mock)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You have to enter a username")

	// Test wrong email
	mock.Add("username", defaultUser.Username)
	mock.Add("email", "wrong_email")
	resp = MemoryRegisterHelper(mock)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You have to enter a valid email address")

	// Test missing and/or non matching passwords
	mock.Set("email", defaultUser.Email)
	mock.Add("password", defaultUser.PwHash)
	mock.Add("password2", "wrong"+defaultUser.PwHash)
	resp = MemoryRegisterHelper(mock)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "The two passwords do not match")

	// Test successful register
	mock.Set("password2", defaultUser.PwHash)
	resp = MemoryRegisterHelper(mock)
	assert.Equal(t, resp.StatusCode, 302, "A successful register should redirrect")
	assert.Equal(t, "/login", resp.Header.Get("Location"))
}
