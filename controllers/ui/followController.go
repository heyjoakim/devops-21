package ui

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
)

// FollowUserHandler handles following another user
func FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	session := GetSession(w, r)
	currentUserID := session.Values["user_id"].(uint)
	params := mux.Vars(r)
	username := params["username"]
	userToFollowID, _ := services.GetUserID(username)

	follower := models.Follower{WhoID: currentUserID, WhomID: userToFollowID}
	err := services.CreateFollower(follower)
	if err != nil {
		log.Fatal(err)
	}

	AddFlash(session, w, r, "You are now following "+username, "Info")
	http.Redirect(w, r, "/"+username, http.StatusFound)
}

// UnfollowUserHandler - relies on a query string
func UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "_cookie")
	loggedInUser := session.Values["user_id"].(uint)
	params := mux.Vars(r)
	username := params["username"]

	if username == "" {
		AddFlash(session, w, r, "No query parameter present", "Warn")
		http.Redirect(w, r, "timeline", http.StatusFound)
		return
	}

	id2, user2Err := services.GetUserID(username)
	if user2Err != nil {
		AddFlash(session, w, r, "User does not exist", "Warn")
		http.Redirect(w, r, "timeline", http.StatusFound)
		return
	}

	err := services.UnfollowUser(loggedInUser, id2)

	if err != nil {
		AddFlash(session, w, r, "Error following user", "Warn")
		fmt.Println("db error: ", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	AddFlash(session, w, r, "You are no longer following "+username, "Info")
	http.Redirect(w, r, "/"+username, http.StatusFound)
}
