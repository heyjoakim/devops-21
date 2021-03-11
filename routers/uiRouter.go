package routers

import (
	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/controllers/ui"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// AddUIRouter creates endpoints for the ui router
func AddUIRouter(router *mux.Router) {
	router.Use(ui.BeforeRequest)
	router.Use(ui.AfterRequest)
	router.HandleFunc("/", ui.TimelineHandler).Methods("GET")
	router.HandleFunc("/{username}/unfollow", ui.UnfollowUserHandler).Methods("GET")
	router.HandleFunc("/{username}/follow", ui.FollowUserHandler).Methods("GET")
	router.HandleFunc("/login", ui.GetLoginHandler).Methods("GET")
	router.HandleFunc("/login", ui.PostLoginHandler).Methods("POST")
	router.HandleFunc("/logout", ui.LogoutHandler).Methods("GET")
	router.HandleFunc("/addMessage", ui.AddMessageHandler).Methods("POST")
	router.HandleFunc("/register", ui.GetRegisterUserHandler).Methods("GET")
	router.HandleFunc("/register", ui.PostRegisterUserHandler).Methods("POST")
	router.HandleFunc("/public", ui.PublicTimelineHandler).Methods("GET")
	router.HandleFunc("/favicon.ico", ui.FaviconHandler).Methods("GET")
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/{username}", ui.UserTimelineHandler).Methods("GET")
}
