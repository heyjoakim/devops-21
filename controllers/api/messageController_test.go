package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

var username = "a"
var email = "a@a.a"
var pass = "a"

// Might be a better way to do this, but follower query cant be exported
// as "no" and "latest" needs to be lowercase...
type LatestMessageRequest struct {
	No     int `json:"no"`
	Latest int `json:"latest"`
}

func TestMemoryApiMessage(t *testing.T) {
	// Register user
	m := models.RegisterRequest{Username: username, Email: email, Password: pass}
	data, _ := json.Marshal(m)
	resp := MemoryRegisterHelper(data)

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)

	// Add message
	c := models.MessageRequest{Content: "Blub!"}
	msg, _ := json.Marshal(c)
	msgResp := MemoryCreateMessageHelper(msg, username)

	assert.Equal(t, msgResp.StatusCode, http.StatusNoContent)

	// Verify that latest was updated
	newResp := MemoryLatestHelper()
	body, _ := ioutil.ReadAll(newResp.Body)

	assert.Contains(t, string(body), `{"latest":2}`)
}

func TestGetLatestUserMessage(t *testing.T) {
	// Prepare query
	var m = LatestMessageRequest{No: 1, Latest: 3}
	query, _ := json.Marshal(m)
	resp := MemoryGetLatestUserMessageHelper(query, username)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Contains(t, string(body), `Blub!`)

	// Verify that latest was updated
	newResp := MemoryLatestHelper()
	newBody, _ := ioutil.ReadAll(newResp.Body)

	assert.Contains(t, string(newBody), `{"latest":3}`)
}

func TestGetLatestMessage(t *testing.T) {
	// Prepare query
	var m = LatestMessageRequest{No: 1, Latest: 4}
	query, _ := json.Marshal(m)
	resp := MemoryGetLatestMessageHelper(query, username)
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, resp.StatusCode, 200)
	assert.Contains(t, string(body), `Blub!`)

	// Verify that latest was updated
	newResp := MemoryLatestHelper()
	newBody, _ := ioutil.ReadAll(newResp.Body)

	assert.Contains(t, string(newBody), `{"latest":4}`)
}
