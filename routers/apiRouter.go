package routers

import (
	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/controllers/api"
)

// AddAPIRoutes creates endpoints for the API router
func AddAPIRoutes(r *mux.Router) {
	r.HandleFunc("/latest", api.GetLatestHandler).Methods("GET")
	r.HandleFunc("/register", api.RegisterHandler).Methods("POST")
	r.HandleFunc("/msgs", api.MessagesHandler).Methods("GET")
	r.HandleFunc("/msgs/{username}", api.GetMessagesFromUserHandler).Methods("GET")
	r.HandleFunc("/msgs/{username}", api.PostMessageHandler).Methods("POST")
	r.HandleFunc("/fllws/{username}", api.GetFollowersHandler).Methods("GET")
	r.HandleFunc("/fllws/{username}", api.FollowHandler).Methods("POST")
}
