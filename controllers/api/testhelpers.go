package api

import (
	"bytes"
	b64 "encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
)

var (
	Credentials        = []string{"simulator", "super_safe!"}
	EncodedCredentials = b64.StdEncoding.EncodeToString([]byte(strings.Join(Credentials, ":")))
	Authentication     = "Basic " + EncodedCredentials
	ContentType        = "Content-Type"
	JSONContent        = "application/json"
)

func SetMuxVars(request *http.Request, username string) *http.Request {
	var vars = map[string]string{
		"username": username,
	}

	return mux.SetURLVars(request, vars)
}

func sendRequest(method string, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add(ContentType, JSONContent)
	req.Header.Add("Authorization", Authentication)
	return req
}

// MemoryRegisterHelepr sends a register user request
func MemoryRegisterHelper(data []byte) *http.Response {
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(data))
	req.Header.Add(ContentType, JSONContent)

	q := req.URL.Query()
	q.Add("latest", "1")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

// MemoryCreateMessageHelper sends a request to create a message
func MemoryCreateMessageHelper(data []byte, username string) *http.Response {
	URI := "/api/msgs/"
	req := sendRequest("POST", URI+username, bytes.NewBuffer(data))

	q := req.URL.Query()
	q.Add("latest", "2")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	req = SetMuxVars(req, username)

	handler := http.HandlerFunc(PostMessageHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

// MemoryGetLatestUserMessageHelper requests to get the latest message from a user
func MemoryGetLatestUserMessageHelper(data []byte, username string) *http.Response {
	URI := "/api/msgs/"
	req := sendRequest("GET", URI+username, bytes.NewBuffer(data))

	q := req.URL.Query()
	q.Add("latest", "3")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	req = SetMuxVars(req, username)

	handler := http.HandlerFunc(GetMessagesFromUserHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

// MemoryGetLatestMessageHelper requests the latest messages
func MemoryGetLatestMessageHelper(data []byte, username string) *http.Response {
	URI := "api/msgs"
	req := sendRequest("GET", URI, bytes.NewBuffer(data))

	q := req.URL.Query()
	q.Add("latest", "4")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	handler := http.HandlerFunc(MessagesHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

// MemoryFollowUserHelper sends a request to follow a user
func MemoryFollowUserHelper(data []byte, username string) *http.Response {
	URI := "/api/fllws/"
	req := sendRequest("POST", URI+username, bytes.NewBuffer(data))

	q := req.URL.Query()
	q.Add("latest", "5")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	req = SetMuxVars(req, username)

	handler := http.HandlerFunc(FollowHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

// MemoryGetFollowUserHelper sends a request to get user followers
func MemoryGetFollowUserHelper(data []byte, username string) *http.Response {
	URI := "/api/fllws/"
	req := sendRequest("GET", URI+username, bytes.NewBuffer(data))

	q := req.URL.Query()
	q.Add("latest", "5")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	req = SetMuxVars(req, username)

	handler := http.HandlerFunc(GetFollowersHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

// MemoryLatestHelper sends a request to get the latest variable
func MemoryLatestHelper() *http.Response {
	URI := "/latest"
	req := sendRequest("GET", URI, nil)

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(GetLatestHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}
