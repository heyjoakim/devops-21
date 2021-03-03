package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var defaultUserData url.Values = url.Values{
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

func MemorySetup() *App {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})

	app := &App{db}
	app.initDb()
	return app
}

// Register a user from a new App
func MemoryRegisterHelper(data url.Values) (*http.Response, *App) {
	app := MemorySetup()
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(app.registerHandler)
	handler.ServeHTTP(w, req)
	return w.Result(), app

}

// Register a user on an existing given app
func RegisterAppHelper(data url.Values, app *App) (*http.Response, *App) {
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(app.registerHandler)
	handler.ServeHTTP(w, req)
	return w.Result(), app

}

// Login user in a existing given app
func MemoryLoginHelper(data url.Values, app *App) (*http.Response, *App) {
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(app.loginHandler)
	handler.ServeHTTP(w, req)
	return w.Result(), app
}

// Register and Login in a new app
func MemoryLoginRegisterHelper(data url.Values) (*http.Response, *App) {
	_, a := MemoryRegisterHelper(data)
	resp, app := MemoryLoginHelper(data, a)
	return resp, app
}

// Add message in a new app
func MemoryAddMessageHelper(data url.Values, registeredUser url.Values) (*http.Response, *App) {
	_, app := MemoryLoginRegisterHelper(registeredUser)
	req, _ := http.NewRequest("POST", "/addMessageHandler", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(app.addMessageHandler)
	handler.ServeHTTP(w, req)
	return w.Result(), app
}

// Add message from user x and user y in a new app
func MemoryTimelineHelper(x url.Values, xdata url.Values, y url.Values, ydata url.Values) (*http.Response, *App) {
	var app *App
	_, app = MemoryAddMessageHelper(xdata, x)
	_, app = RegisterAppHelper(y, app)
	_, app = MemoryLoginHelper(y, app)

	req, _ := http.NewRequest("POST", "/addMessageHandler", strings.NewReader(ydata.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(app.addMessageHandler)
	handler.ServeHTTP(w, req)
	return w.Result(), app
}

func TestMemoryRegister(t *testing.T) {
	var resp *http.Response
	var body []byte
	mock := url.Values{}

	// Test missing username
	resp, _ = MemoryRegisterHelper(mock)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You have to enter a username")

	// Test wrong email
	mock.Add("username", defaultUser.Username)
	mock.Add("email", "wrong_email")
	resp, _ = MemoryRegisterHelper(mock)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You have to enter a valid email address")

	// Test missing and/or non matching passwords
	mock.Set("email", defaultUser.Email)
	mock.Add("password", defaultUser.PwHash)
	mock.Add("password2", "wrong"+defaultUser.PwHash)
	resp, _ = MemoryRegisterHelper(mock)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "The two passwords do not match")

	// Test successful register
	mock.Set("password2", defaultUser.PwHash)
	resp, _ = MemoryRegisterHelper(mock)
	assert.Equal(t, resp.StatusCode, 302, "A successful register should redirrect")
	assert.Equal(t, "/login", resp.Header.Get("Location"))

}

func TestMemoryLoginHelper(t *testing.T) {
	var resp *http.Response
	var body []byte
	mock := url.Values{}

	// Need to register a user to test error message
	_, app := MemoryRegisterHelper(defaultUserData)

	// Test wrong username
	mock.Add("username", "wrong_username")
	resp, _ = MemoryLoginHelper(mock, app)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid username")

	// Test missing password
	mock.Set("username", defaultUser.Username)
	mock.Add("password", "wrong_password")
	resp, _ = MemoryLoginHelper(mock, app)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid password")

	// Test successful login
	mock.Set("password", defaultUser.PwHash)
	resp, _ = MemoryLoginHelper(mock, app)
	assert.Equal(t, resp.StatusCode, 302, "A successful login should redirrect")
}

func TestMemoryLogout(t *testing.T) {
	_, app := MemoryLoginRegisterHelper(defaultUserData)

	// Setup cookie
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	session, _ := store.Get(req, "_cookie")
	session.Values["user_id"] = defaultUser.UserID
	session.Save(req, w)
	cookie := session.Values["user_id"]

	// Assert that a cookie is actually set
	assert.Equal(t, cookie, defaultUser.UserID)

	// Serve request
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler := http.HandlerFunc(app.logoutHandler)
	handler.ServeHTTP(w, req)
	emptyCookie := session.Values["user_id"]

	logoutResponse := w.Result()

	// Assert that the cookie is now empty and redirrect
	assert.NotEqual(t, emptyCookie, defaultUser.UserID)
	assert.Equal(t, logoutResponse.StatusCode, 302)
}

func TestMemoryAddMessage(t *testing.T) {
	var resp *http.Response

	// Expected message
	var msg string = "Test message personal page"
	var msgData url.Values = url.Values{
		"text":  {msg},
		"token": {defaultUser.Username},
	}

	// Add message
	resp, app := MemoryAddMessageHelper(msgData, defaultUserData)

	// Assert that adding a message redirects
	assert.Equal(t, resp.StatusCode, 302)

	// Beigin new request to check the users page for the added message
	req, _ := http.NewRequest("GET", "/username", nil)
	w := httptest.NewRecorder()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set URL vars to be retrieved by mux.Vars
	req = mux.SetURLVars(req, map[string]string{"username": defaultUser.Username})

	// Set session values
	session, _ := store.Get(req, "_cookie")
	session.Values["user_id"] = defaultUser.UserID
	session.Values["username"] = defaultUser.Username
	session.Save(req, w)

	// Handle request
	handler := http.HandlerFunc(app.userTimelineHandler)
	handler.ServeHTTP(w, req)
	checkResp := w.Result()
	body, _ := ioutil.ReadAll(checkResp.Body)

	// Assert that new message is added to the personal page
	assert.Contains(t, string(body), msg)

}

func TestMemoryTimeline(t *testing.T) {

	// Expected message
	var msg string = "Test message on timeline"
	var msgData url.Values = url.Values{
		"text":  {msg},
		"token": {defaultUser.Username},
	}

	// Add message
	_, app := MemoryAddMessageHelper(msgData, defaultUserData)

	// Begin new request to check the public page for message
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Set session values
	session, _ := store.Get(req, "_cookie")
	session.Values["user_id"] = defaultUser.UserID
	session.Values["username"] = defaultUser.Username
	session.Save(req, w)

	handler := http.HandlerFunc(app.publicTimelineHandler)
	handler.ServeHTTP(w, req)
	checkResp := w.Result()
	body, _ := ioutil.ReadAll(checkResp.Body)

	// Assert that new message is added to the page
	assert.Contains(t, string(body), msg)

}

func TestMemoryFollow(t *testing.T) {
	// Setup two mock users
	foo := &models.User{Username: "Foo", Email: "foo@baz.com", PwHash: "off"}
	bar := &models.User{Username: "Bar", Email: "bar@baz.com", PwHash: "rab"}
	fooData := url.Values{"username": {foo.Username}, "email": {foo.Email}, "password": {foo.PwHash}, "password2": {foo.PwHash}}
	barData := url.Values{"username": {bar.Username}, "email": {bar.Email}, "password": {bar.PwHash}, "password2": {bar.PwHash}}

	_, app := MemoryTimelineHelper(
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
	session.Values["user_id"], _ = app.getUserID(bar.Username)
	session.Values["username"] = bar.Username
	session.Save(req, w)

	// Set URL vars to be retrieved by mux.Vars
	// Expected params["username"] = foo
	req = mux.SetURLVars(req, map[string]string{"username": foo.Username})

	handler := http.HandlerFunc(app.followUserHandler)
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
	session.Values["user_id"], _ = app.getUserID(bar.Username)
	session.Values["username"] = bar.Username
	session.Save(newReq, newW)

	// Handle request
	newHandler := http.HandlerFunc(app.userTimelineHandler)
	newHandler.ServeHTTP(newW, newReq)
	newResp := newW.Result()
	newBody, _ := ioutil.ReadAll(newResp.Body)

	// Assert that new message is added to the personal page
	assert.Contains(t, string(newBody), "Foo test message")
	assert.Contains(t, string(newBody), "Bar test message")

}
