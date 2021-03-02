package api

import (
	"net/http"
	"strconv"
)

var (
	latest = 0
)

func updateLatest(r *http.Request) {
	tryLatestQuery := r.URL.Query().Get("latest")

	if tryLatestQuery == "" {
		latest = -1
	} else {
		tryLatest, _ := strconv.Atoi(tryLatestQuery)
		latest = tryLatest
	}
}