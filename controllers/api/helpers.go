package api

import (
	"net/http"
	"strconv"

	"github.com/heyjoakim/devops-21/services"
)

var d = services.GetDbInstance()

func updateLatest(r *http.Request) {
	tryLatestQuery := r.URL.Query().Get("latest")

	if tryLatestQuery != "" {
		tryLatest, _ := strconv.Atoi(tryLatestQuery)

		services.UpdateLatest(tryLatest)
	}
}
