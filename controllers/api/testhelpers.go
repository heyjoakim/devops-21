package api

import (
	"bytes"
	b64 "encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
)

var (
	CREDENTIALS         = []string{"simulator", "super_safe!"}
	ENCODED_CREDENTILAS = b64.StdEncoding.EncodeToString([]byte(strings.Join(CREDENTIALS, ":")))
	JSONCONTENT         = "application/json"
)

func MemoryRegisterHelper(data []byte) *http.Response {
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(data))
	req.Header.Add("Content-Type", JSONCONTENT)

	q := req.URL.Query()
	q.Add("latest", "1")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

func MemoryCreateMessageHelper(data []byte, username string) *http.Response {
	URI := "/api/msgs/"
	req, _ := http.NewRequest("POST", URI+username, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", JSONCONTENT)
	req.Header.Add("Authorization", "Basic "+ENCODED_CREDENTILAS)

	q := req.URL.Query()
	q.Add("latest", "2")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	// Set MUX vars
	var vars = map[string]string{
		"username": username,
	}

	req = mux.SetURLVars(req, vars)

	handler := http.HandlerFunc(PostMessageHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

func MemoryGetLatestUserMessageHelper(data []byte, username string) *http.Response {
	URI := "/api/msgs/"

	req, _ := http.NewRequest("GET", URI+username, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", JSONCONTENT)
	req.Header.Add("Authorization", "Basic "+ENCODED_CREDENTILAS)

	q := req.URL.Query()
	q.Add("latest", "3")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	// Set MUX vars
	var vars = map[string]string{
		"username": username,
	}

	req = mux.SetURLVars(req, vars)

	handler := http.HandlerFunc(GetMessagesFromUserHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

func MemoryGetLatestMessageHelper(data []byte, username string) *http.Response {
	req, _ := http.NewRequest("GET", "api/msgs", bytes.NewBuffer(data))
	req.Header.Add("Content-Type", JSONCONTENT)
	req.Header.Add("Authorization", "Basic "+ENCODED_CREDENTILAS)

	q := req.URL.Query()
	q.Add("latest", "4")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	handler := http.HandlerFunc(MessagesHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

func MemoryFollowUserHelper(data []byte, username string) *http.Response {
	URI := "/api/fllws/"
	req, _ := http.NewRequest("POST", URI+username, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", JSONCONTENT)
	req.Header.Add("Authorization", "Basic "+ENCODED_CREDENTILAS)

	q := req.URL.Query()
	q.Add("latest", "5")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	// Set MUX vars
	var vars = map[string]string{
		"username": username,
	}

	req = mux.SetURLVars(req, vars)

	handler := http.HandlerFunc(FollowHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

func MemoryGetFollowUserHelper(data []byte, username string) *http.Response {
	URI := "/api/fllws/"
	req, _ := http.NewRequest("GET", URI+username, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", JSONCONTENT)
	req.Header.Add("Authorization", "Basic "+ENCODED_CREDENTILAS)

	q := req.URL.Query()
	q.Add("latest", "5")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	// Set MUX vars
	var vars = map[string]string{
		"username": username,
	}

	req = mux.SetURLVars(req, vars)

	handler := http.HandlerFunc(GetFollowersHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

func MemoryLatestHelper() *http.Response {
	req, _ := http.NewRequest("GET", "/latest", nil)
	req.Header.Add("Content-Type", JSONCONTENT)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(GetLatestHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}
