package api

import (
	"fmt"
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

func RegisterEndpoint(name string) *prometheus.Timer {
	var (
		bucketStart = 0.01
		bucketWidth = 0.05
		bucketCount = 10
	)

	hist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("http_request_%s_duration_seconds", name),
			Help:    fmt.Sprintf("http_request_%s_duration_seconds", name),
			Buckets: prometheus.LinearBuckets(bucketStart, bucketWidth, bucketCount),
		},
		[]string{"status"},
	)
	var status string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		hist.WithLabelValues(status).Observe(v)
	}))
	prometheus.MustRegister(hist)
	return timer
}
