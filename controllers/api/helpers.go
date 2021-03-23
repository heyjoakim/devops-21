package api

import (
	"net/http"
	"strconv"

	"github.com/heyjoakim/devops-21/services"
	"github.com/prometheus/client_golang/prometheus"
)

func updateLatest(r *http.Request) {
	tryLatestQuery := r.URL.Query().Get("latest")

	if tryLatestQuery != "" {
		tryLatest, _ := strconv.Atoi(tryLatestQuery)

		services.UpdateLatest(tryLatest)
	}
}

func createEndpointTimer(hist *prometheus.HistogramVec) *prometheus.Timer {
	var status string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		hist.WithLabelValues(status).Observe(v)
	}))
	return timer
}
