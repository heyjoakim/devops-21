package ui

import (
	"fmt"
	"net/http"

	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
)

// GetRegisterUserHandler returns the register page..
func GetRegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	session := GetSession(w, r)
	if ok := session.Values["user_id"] != nil; ok {
		routeName := fmt.Sprintf("/%s", session.Values["username"])
		http.Redirect(w, r, routeName, http.StatusFound)
	}

	data := models.PageData{
		"username": session.Values["username"],
	}
	redirectToRegister(w, data)
}

// PostRegisterUserHandler handles user signup requests.
func PostRegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	session := GetSession(w, r)
	if ok := session.Values["user_id"] != nil; ok {
		routeName := fmt.Sprintf("/%s", session.Values["username"])
		http.Redirect(w, r, routeName, http.StatusFound)
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	user := models.UserCreateRequest{
		Username:  username,
		Email:     email,
		Password:  password,
		Password2: password2,
	}
	err := services.CreateUser(user)

	if err != nil {
		data := models.PageData{
			"error": err,
		}
		redirectToRegister(w, data)
		return
	}

	AddFlash(session, w, r, "You are now registered ?", username)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func redirectToRegister(w http.ResponseWriter, data models.PageData) {
	tmpl := LoadTemplate(RegisterPath)
	_ = tmpl.Execute(w, data)
}
