package ui

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

var logoutUserData url.Values = url.Values{
	"username":  {"Tom"},
	"email":     {"tom@go.com"},
	"password":  {"secret"},
	"password2": {"secret"},
}

var logoutUser = &models.User{
	Username: "Tom",
	Email:    "tom@go.com",
	PwHash:   "secret",
}

func TestMemoryLogout(t *testing.T) {
	MemoryLoginRegisterHelper(logoutUserData)

	// Setup cookie
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	session, _ := store.Get(req, "_cookie")
	session.Values["user_id"] = logoutUser.UserID
	_ = session.Save(req, w)
	cookie := session.Values["user_id"]

	// Assert that a cookie is actually set
	assert.Equal(t, cookie, logoutUser.UserID)

	// Serve request
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler := http.HandlerFunc(LogoutHandler)
	handler.ServeHTTP(w, req)
	emptyCookie := session.Values["user_id"]

	logoutResponse := w.Result()

	// Assert that the cookie is now empty and redirrect
	assert.NotEqual(t, emptyCookie, logoutUser.UserID)
	assert.Equal(t, logoutResponse.StatusCode, 302)
}
