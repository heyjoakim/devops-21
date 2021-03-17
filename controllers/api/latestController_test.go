package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

func TestMemoryApiLatest(t *testing.T) {
	// Register a user
	m := models.RegisterRequest{Username: "test", Email: "test@test", Password: "foo"}
	data, _ := json.Marshal(m)
	resp := MemoryRegisterHelper(data)

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)

	// verify latest variable
	newResp := MemoryLatestHelper()
	body, _ := ioutil.ReadAll(newResp.Body)

	assert.Contains(t, string(body), `{"latest":1}`)
}
