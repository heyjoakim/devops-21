package api

import (
	"github.com/heyjoakim/devops-21/helpers"
	"net/http"
)


// GetLatest godoc
// @Summary Get the latest x
// @Description Get the latest x
// @Produce  json
// @Success 200 {object} interface{}
// @Router /api/latest [get]
func GetLatestHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"latest": latest,
	}

	jsonData, err := helpers.Serialize(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
