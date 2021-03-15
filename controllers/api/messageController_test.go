package api

import (
	"testing"
)

func TestMemoryApiMessage(t *testing.T) {

	// username := "a"

	// data, _ := json.Marshal(models.MessageRequest{
	// 	Content: "Blub!",
	// })

	// req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(data))
	// req.Header.Set("Content-Type", "application/json")

	// q := req.URL.Query()
	// q.Add("latest", "2")
	// req.URL.RawQuery = q.Encode()

	// w := httptest.NewRecorder()
	// handler := http.HandlerFunc(RegisterHandler)
	// handler.ServeHTTP(w, req)
	// resp := w.Result()

	// assert.Equal(t, resp.StatusCode, http.StatusNoContent)

	// newReq, _ := http.NewRequest("GET", "/latest", nil)
	// newW := httptest.NewRecorder()
	// newHandler := http.HandlerFunc(GetLatestHandler)
	// newHandler.ServeHTTP(newW, newReq)
	// newResp := newW.Result()
	// body, _ := ioutil.ReadAll(newResp.Body)

}
