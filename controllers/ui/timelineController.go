package ui

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	log "github.com/sirupsen/logrus"
)

// TimelineHandler a users timeline or if no user is logged in it will
// redirect to the public timeline.  This timeline shows the user's
// messages as well as all the messages of followed users.
func TimelineHandler(w http.ResponseWriter, r *http.Request) {
	session := GetSession(w, r)
	if session.Values["user_id"] != nil {
		routeName := fmt.Sprintf("/%s", session.Values["username"])
		http.Redirect(w, r, routeName, http.StatusFound)
		return
	}

	http.Redirect(w, r, "/public", http.StatusFound)
}

// PublicTimelineHandler shows the public timeline
func PublicTimelineHandler(w http.ResponseWriter, r *http.Request) {
	var results = services.GetPublicMessages(PerPage)
	var messages []models.MessageViewModel
	for _, result := range results {
		message := models.MessageViewModel{
			Email:   helpers.GetGravatarURL(result.Email, 48),
			User:    result.Username,
			Content: result.Content,
			PubDate: helpers.FormatDatetime(result.PubDate),
		}
		messages = append(messages, message)
	}

	session := GetSession(w, r)
	username := session.Values["username"]

	data := models.PageData{
		"username": username,
		"messages": messages,
		"msgCount": len(messages),
	}

	tmpl := LoadTemplate(TimelinePath)
	_ = tmpl.Execute(w, data)
}

// UserTimelineHandler shows the posts from one user
func UserTimelineHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	profileUsername := params["username"]

	profileUserID, err := services.GetUserID(profileUsername)
	if err != nil {
		errorCode := 404
		w.WriteHeader(errorCode)
		log.WithField("err", err).Error("UserTimelineHandler error")
		return
	}

	session := GetSession(w, r)
	sessionUserID := session.Values["user_id"]
	data := models.PageData{"followed": false}

	if sessionUserID != nil {
		if services.IsUserFollower(sessionUserID.(uint), profileUserID) {
			data["followed"] = true
		}
	}

	messages := getPostsForUser(profileUserID)

	var msgS []models.MessageViewModel
	if sessionUserID != nil {
		if sessionUserID == profileUserID {
			followlist := getFollowedUsers(sessionUserID.(uint))
			for _, v := range followlist {
				msgS = append(getPostsForUser(v), msgS...)
			}
		}
	}

	messages = append(msgS, messages...)
	data["messages"] = messages
	data["title"] = fmt.Sprintf("%s's Timeline", profileUsername)
	data["profileOwner"] = profileUsername

	if session.Values["username"] == profileUsername {
		data["ownProfile"] = true
	}

	data["msgCount"] = len(messages)
	data["username"] = session.Values["username"]
	data["MsgInfo"] = session.Flashes("Info")
	data["MsgWarn"] = session.Flashes("Warn")
	_ = session.Save(r, w)

	tmpl, err := template.ParseFiles(TimelinePath, LayoutPath)
	if err != nil {
		log.Error(err)
	}
	_ = tmpl.Execute(w, data)
}

func getPostsForUser(id uint) []models.MessageViewModel {
	var results = services.GetMessagesForUser(PerPage, id)

	var messages []models.MessageViewModel
	for _, result := range results {
		message := models.MessageViewModel{
			Email:   helpers.GetGravatarURL(result.Email, 48),
			User:    result.Username,
			Content: result.Content,
			PubDate: helpers.FormatDatetime(result.PubDate),
		}
		messages = append(messages, message)
	}

	return messages
}

// get ID's of all users that are followed by some user
func getFollowedUsers(id uint) []uint {
	var followers = services.GetUsersFollowedBy(id)

	var followlist []uint
	for _, follower := range followers {
		followlist = append(followlist, follower.WhomID)
	}

	return followlist
}
