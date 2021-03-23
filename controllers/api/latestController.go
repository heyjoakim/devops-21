package api

import (
	"net/http"

	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/metrics"
	"github.com/heyjoakim/devops-21/services"
)

// GetLatestHandler godoc
// @Summary Get the latest x
// @Description Get the latest x
// @Produce  json
// @Success 200 {object} interface{}
// @Router /api/latest [get]
func GetLatestHandler(w http.ResponseWriter, r *http.Request) {
	hist := metrics.GetHistogramVec("get_api_latest")
	if hist != nil {
		timer := createEndpointTimer(hist)
		defer timer.ObserveDuration()
	}

	data := map[string]interface{}{
		"latest": services.GetLatest(),
	}

	jsonData, err := helpers.Serialize(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonData)
}
