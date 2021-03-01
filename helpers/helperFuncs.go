package helpers

import (
	"encoding/json"
	"net/http"
)

func GetQueryParameter(r *http.Request, name string) string {
	queryValue := r.URL.Query().Get(name)
	return queryValue
}

func Serialize(input interface{}) ([]byte, error) {
	js, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func NotReqFromSimulator(r *http.Request) []byte {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		data := map[string]interface{}{
			"status":    http.StatusForbidden,
			"error_msg": "You are not authorized to use this resource!",
		}

		jsonData, _ := Serialize(data)
		return jsonData
	}
	return nil
}
