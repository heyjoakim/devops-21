package routers

import (
	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/controllers/ui"
)

// AddUIRouter creates endpoints for the ui router
func AddUIRouter(router *mux.Router) {
	router.Use(ui.BeforeRequest)
	router.Use(ui.AfterRequest)
	router.HandleFunc("/", ui.TimelineHandler)
	router.HandleFunc("/{username}/unfollow", ui.UnfollowUserHandler).Methods("GET")
	router.HandleFunc("/{username}/follow", ui.FollowUserHandler).Methods("GET")
	router.HandleFunc("/login", ui.GetLoginHandler).Methods("GET")
	router.HandleFunc("/login", ui.PostLoginHandler).Methods("POST")
	router.HandleFunc("/logout", ui.LogoutHandler)
	router.HandleFunc("/addMessage", ui.AddMessageHandler).Methods("POST")
	router.HandleFunc("/register", ui.GetRegisterUserHandler).Methods("GET")
	router.HandleFunc("/register", ui.PostRegisterUserHandler).Methods("POST")
	router.HandleFunc("/public", ui.PublicTimelineHandler)
	router.HandleFunc("/favicon.ico", ui.FaviconHandler)
	router.HandleFunc("/{username}", ui.UserTimelineHandler)
}
