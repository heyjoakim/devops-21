package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

var userA = "fa"
var emailA = "fa@fa"
var passA = "fa"

var userB = "fb"
var emailB = "fb@fb"
var passB = "fb"

// Might be a better way to do this, but follower query cant be exported
// as "no" and "latest" needs to be lowercase...
type GetFollowerQuery struct {
	No     int `json:"no"`
	Latest int `json:"latest"`
}

func TestMemoryApiFollow(t *testing.T) {
	// Register users
	ma := models.RegisterRequest{Username: userA, Email: emailA, Password: passA}
	mb := models.RegisterRequest{Username: userB, Email: emailB, Password: passB}

	dataA, _ := json.Marshal(ma)
	dataB, _ := json.Marshal(mb)

	MemoryRegisterHelper(dataA)
	MemoryRegisterHelper(dataB)

	// Send follow request
	m := models.FollowRequest{Follow: userB}
	followData, _ := json.Marshal(m)
	resp := MemoryFollowUserHelper(followData, userA)

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)

	// Query get follower
	query := GetFollowerQuery{No: 1, Latest: 6}
	dataQuery, _ := json.Marshal(query)
	newResp := MemoryGetFollowUserHelper(dataQuery, userA)
	body, _ := ioutil.ReadAll(newResp.Body)

	assert.Contains(t, string(body), "fb")
}

func TestMemoryApiUnfollow(t *testing.T) {
	// Send unfollow request
	m := models.FollowRequest{Unfollow: userB}
	data, _ := json.Marshal(m)
	resp := MemoryFollowUserHelper(data, userA)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.NotContains(t, string(body), "fb")
}
