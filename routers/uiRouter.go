package routers

import (
	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/controllers"
	"github.com/heyjoakim/devops-21/controllers/ui"
)

// AddUIRouter creates endpoints for the ui router
func AddUIRouter(router *mux.Router) {
	router.Use(controllers.BeforeRequest)
	router.Use(controllers.AfterRequest)
	router.HandleFunc("/", ui.TimelineHandler)
	router.HandleFunc("/{username}/unfollow", ui.UnfollowUserHandler).Methods("GET")
	router.HandleFunc("/{username}/follow", ui.FollowUserHandler).Methods("GET")
	router.HandleFunc("/login", ui.GetLoginHandler).Methods("GET")
	router.HandleFunc("/login", ui.PostLoginHandler).Methods("POST")
	router.HandleFunc("/logout", controllers.LogoutHandler)
	router.HandleFunc("/addMessage", controllers.AddMessageHandler).Methods("GET", "POST")
	router.HandleFunc("/register", controllers.RegisterUserUiHandler).Methods("GET", "POST")
	router.HandleFunc("/public", ui.PublicTimelineHandler)
	router.HandleFunc("/favicon.ico", controllers.FaviconHandler)
	router.HandleFunc("/{username}", ui.UserTimelineHandler)
}
