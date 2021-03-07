package api

import (
	"net/http"
	"strconv"

	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
)

var (
	latest = 0
)
var d = services.GetDbInstance()

func updateLatest(r *http.Request) {
	tryLatestQuery := r.URL.Query().Get("latest")

	if tryLatestQuery == "" {
		latest = -1
	} else {
		tryLatest, _ := strconv.Atoi(tryLatestQuery)

		latest = tryLatest
	}

	var c models.Config

	er := d.DB.First(&c, "key = ?", "latest").Error
	if er != nil {

		c.ID = 0
		c.Key = "latest"
		c.Value = strconv.Itoa(latest)
		d.DB.Create(&c)
	} else {

		d.DB.Model(&models.Config{}).Where("key = ?", "latest").Update("Value", latest)
	}
}
