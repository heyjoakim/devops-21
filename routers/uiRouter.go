package routers

import (
	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/controllers"
)

func AddUiRouter(router *mux.Router) {
	router.Use(controllers.BeforeRequest)
	router.Use(controllers.AfterRequest)
	router.HandleFunc("/", controllers.TimelineHandler)
	router.HandleFunc("/{username}/unfollow", controllers.UnfollowUserHandler)
	router.HandleFunc("/{username}/follow", controllers.FollowUserHandler)
	router.HandleFunc("/login", controllers.LoginHandler).Methods("GET", "POST")
	router.HandleFunc("/logout", controllers.LogoutHandler)
	router.HandleFunc("/addMessage", controllers.AddMessageHandler).Methods("GET", "POST")
	router.HandleFunc("/register", controllers.RegisterUserUiHandler).Methods("GET", "POST")
	router.HandleFunc("/public", controllers.PublicTimelineHandler)
	router.HandleFunc("/favicon.ico", controllers.FaviconHandler)
	router.HandleFunc("/{username}", controllers.UserTimelineHandler)
}
