package ui

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// MemoryRegisterHelper registers a user from a new App
func MemoryRegisterHelper(data url.Values) *http.Response {
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(PostRegisterUserHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

// RegisterAppHelper registers a user on an existing given app
func RegisterAppHelper(data url.Values) *http.Response {
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(PostRegisterUserHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

// MemoryLoginHelper logins user in a existing given app
func MemoryLoginHelper(data url.Values) *http.Response {
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(PostLoginHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

// MemoryLoginRegisterHelper registers and Login in a new app
func MemoryLoginRegisterHelper(data url.Values) *http.Response {
	MemoryRegisterHelper(data)
	resp := MemoryLoginHelper(data)
	return resp
}

// MemoryAddMessageHelper adds message in a new app
func MemoryAddMessageHelper(data url.Values, registeredUser url.Values) *http.Response {
	MemoryLoginRegisterHelper(registeredUser)
	req, _ := http.NewRequest("POST", "/addMessageHandler", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(AddMessageHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}

// MemoryTimelineHelper adds message from user x and user y in a new app
func MemoryTimelineHelper(x url.Values, xdata url.Values, y url.Values, ydata url.Values) *http.Response {
	MemoryAddMessageHelper(xdata, x)
	RegisterAppHelper(y)
	MemoryLoginHelper(y)

	req, _ := http.NewRequest("POST", "/addMessageHandler", strings.NewReader(ydata.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(AddMessageHandler)
	handler.ServeHTTP(w, req)
	return w.Result()
}
