package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

var username_a = "fa"
var email_a = "fa@fa"
var password_a = "fa"

var username_b = "fb"
var email_b = "fb@fb"
var password_b = "fb"

// Congratulations to golang on defining such a stupid way of exporting
type GetFollowerQuery struct {
	no     int `json:"no"`
	latest int `json:"latest"`
}

func TestMemoryApiFollow(t *testing.T) {
	m_a := models.RegisterRequest{Username: username_a, Email: email_a, Password: password_a}
	m_b := models.RegisterRequest{Username: username_b, Email: email_b, Password: password_b}

	data_a, _ := json.Marshal(m_a)
	data_b, _ := json.Marshal(m_b)

	MemoryRegisterHelper(data_a)
	MemoryRegisterHelper(data_b)
	m := models.FollowRequest{Follow: username_b}
	data_f, _ := json.Marshal(m)
	resp := MemoryFollowUserHelper(data_f, username_a)

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)

	query := GetFollowerQuery{no: 1, latest: 6}
	data_query, _ := json.Marshal(query)
	newResp := MemoryGetFollowUserHelper(data_query, username_a)
	body, _ := ioutil.ReadAll(newResp.Body)

	assert.Contains(t, string(body), "fb")
}

func TestMemoryApiUnfollow(t *testing.T) {
	m := models.FollowRequest{Unfollow: username_b}
	data, _ := json.Marshal(m)
	resp := MemoryFollowUserHelper(data, username_a)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.NotContains(t, string(body), "fb")

}
