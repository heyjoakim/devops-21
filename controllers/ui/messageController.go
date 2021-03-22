package ui

import (
	"net/http"
	"time"

	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	log "github.com/sirupsen/logrus"
)

// AddMessageHandler adds a new message to the database.
func AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := services.GetUserID(r.FormValue("token"))
	message := models.Message{AuthorID: userID, Text: r.FormValue("text"), PubDate: time.Now().Unix(), Flagged: 0}
	err := services.CreateMessage(message)
	if err != nil {
		log.WithField("err", err).Error("AddMessageHandler: DB err")
		return
	}

	http.Redirect(w, r, "/"+r.FormValue("token"), http.StatusFound)
}
