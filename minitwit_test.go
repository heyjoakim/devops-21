package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var data url.Values = url.Values{
	"username":  {"Richard"},
	"email":     {"richard@stallman.org"},
	"password":  {"secret"},
	"password2": {"secret"}}

var usr = &User{
	userID:   244,
	username: "Richard",
	email:    "richard@stallman.org",
	pwHash:   "secret",
}

func Setup() (*sql.DB, sqlmock.Sqlmock) {
	tdb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to initialize mock db with error '%s'", err)
	}

	return tdb, mock
}

func TestGetUserID(t *testing.T) {
	tdb, mock := Setup()
	app := &App{tdb}
	defer app.db.Close()

	// Check for expected query
	query := regexp.QuoteMeta(`select user_id from user where username = ?`)
	rows := sqlmock.NewRows([]string{"user_id"}).AddRow(usr.userID)
	mock.ExpectQuery(query).WithArgs(usr.username).WillReturnRows(rows)

	// Get ID of user
	ID, err := app.getUserID(usr.username)

	// Assert that userID gotten from mock db is same as returned
	assert.NoError(t, err)
	assert.Equal(t, usr.userID, ID, "IDs should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("All expectations were not met: %s", err)
	}
}

func TestRegisterHandler(t *testing.T) {
	tdb, mock := Setup()
	app := &App{tdb}
	defer app.db.Close()

	// Create new request and record it
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Check for expected prepare statement
	query := regexp.QuoteMeta(`insert into user (username, email, pw_hash) values(?,?,?)`)
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec() // TODO: find out if we can make expected result for a even more robust test

	// Handle request
	handler := http.HandlerFunc(app.registerHandler)
	handler.ServeHTTP(w, req)
	resp := w.Result()

	// Assert that we get a reddirect (statuscode 302)
	assert.Equal(t, resp.StatusCode, 302)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("All expectations were not met: %s", err)
	}
}

func TestLoginHandler(t *testing.T) {
	tdb, mock := Setup()
	app := &App{tdb}
	defer app.db.Close()

	hash, _ := bcrypt.GenerateFromPassword([]byte(usr.pwHash), bcrypt.DefaultCost)

	// Create new request and record it
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Check for expected query statement
	query := regexp.QuoteMeta(`select * from user where username = ?`)
	rows := sqlmock.NewRows([]string{"user_id", "username", "email", "pw_hash"}).AddRow(usr.userID, usr.username, usr.email, hash)
	mock.ExpectQuery(query).WithArgs(usr.username).WillReturnRows(rows)

	// Handle request
	handler := http.HandlerFunc(app.loginHandler)
	handler.ServeHTTP(w, req)
	resp := w.Result()

	// Assert that we get a reddirect (statuscode 302)
	assert.Equal(t, resp.StatusCode, 302)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("All expectations were not met: %s", err)
	}

}

func TestLogoutHandler(t *testing.T) {
	tdb, mock := Setup()
	app := &App{tdb}
	defer app.db.Close()

	// Setup cookie
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	session, err := store.Get(req, "_cookie")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	session.Values["user_id"] = usr.userID
	session.Save(req, w)
	cookie := session.Values["user_id"]

	// Assert that a cookie is actually set
	assert.Equal(t, cookie, usr.userID)

	// Serve request
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler := http.HandlerFunc(app.logoutHandler)
	handler.ServeHTTP(w, req)
	emptyCookie := session.Values["user_id"]

	resp := w.Result()

	// Assert that the cookie is now empty and redirrect
	assert.NotEqual(t, emptyCookie, usr.userID)
	assert.Equal(t, resp.StatusCode, 302)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("All expectations were not met: %s", err)
	}

}
