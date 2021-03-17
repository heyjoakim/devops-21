package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/heyjoakim/devops-21/models"
	"github.com/stretchr/testify/assert"
)

func TestMemoryApiRegister(t *testing.T) {
	// Some request updating latest
	m := models.RegisterRequest{Username: "foo", Email: "foo@foo", Password: "foo"}
	data, _ := json.Marshal(m)
	resp := MemoryRegisterHelper(data)

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)

	// Verify that latest was updated
	newResp := MemoryLatestHelper()
	body, _ := ioutil.ReadAll(newResp.Body)

	assert.Contains(t, string(body), `{"latest":1}`)
}
