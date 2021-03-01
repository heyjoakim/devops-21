package routers

import (
	"github.com/gorilla/mux"
	"github.com/heyjoakim/devops-21/controllers"
)

func AddAPIRoutes(r *mux.Router) {
	r.HandleFunc("/latest", controllers.GetLatestHandler)
	r.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")
	r.HandleFunc("/msgs", controllers.MessagesHandler)
	r.HandleFunc("/msgs/{username}", controllers.MessagesPerUserHandler).Methods("GET", "POST")
	r.HandleFunc("/fllws/{username}", controllers.FollowHandler).Methods("GET", "POST")

}
