package api

import (
	"net/http"
	"strconv"

	"github.com/heyjoakim/devops-21/services"
)

func updateLatest(r *http.Request) {
	tryLatestQuery := r.URL.Query().Get("latest")

	if tryLatestQuery != "" {
		tryLatest, _ := strconv.Atoi(tryLatestQuery)

		services.UpdateLatest(tryLatest)
	}
}
