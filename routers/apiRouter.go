package routers

import (
	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/controllers"
	"github.com/heyjoakim/devops-21/controllers/api"
)

// AddAPIRoutes creates endpoints for the API router
func AddAPIRoutes(r *mux.Router) {
	r.HandleFunc("/latest", controllers.GetLatestHandler)
	r.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")
	r.HandleFunc("/msgs", controllers.MessagesHandler)
	r.HandleFunc("/msgs/{username}", controllers.MessagesPerUserHandler).Methods("GET", "POST")
	r.HandleFunc("/fllws/{username}", api.GetFollowersHandler).Methods("GET")
	r.HandleFunc("/fllws/{username}", api.FollowHandler).Methods("POST")
}
