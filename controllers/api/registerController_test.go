package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

func TestMemoryApiRegister(t *testing.T) {

	data, _ := json.Marshal(models.RegisterRequest{
		Username: "test",
		Email:    "test@test",
		Password: "foo",
	})

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("latest", "1")
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.Body)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)

	newReq, _ := http.NewRequest("GET", "/latest", nil)
	newW := httptest.NewRecorder()
	newHandler := http.HandlerFunc(GetLatestHandler)
	newHandler.ServeHTTP(newW, newReq)
	newResp := newW.Result()
	body, _ := ioutil.ReadAll(newResp.Body)
	fmt.Println(string(body))

	assert.Contains(t, string(body), `{"latest":1}`)

}
