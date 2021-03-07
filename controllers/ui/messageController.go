package ui

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
)

// AddMessageHandler adds a new message to the database.
func AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := services.GetUserID(r.FormValue("token"))
	message := models.Message{AuthorID: userID, Text: r.FormValue("text"), PubDate: time.Now().Unix(), Flagged: 0}
	err := services.CreateMessage(message)
	if err != nil {
		fmt.Println("database error: ", err)
		log.Fatal(err)
		return
	}

	http.Redirect(w, r, "/"+r.FormValue("token"), http.StatusFound)
}
