package ui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// GetRegisterUserHandler returns the register page.
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

	var registerError string
	if len(r.FormValue("username")) == 0 {
		registerError = "You have to enter a username"
		} else if len(r.FormValue("email")) == 0 || !strings.Contains(r.FormValue("email"), "@") {
			registerError = "You have to enter a valid email address"
			} else if len(r.FormValue("password")) == 0 {
				registerError = "You have to enter a password"
				} else if r.FormValue("password") != r.FormValue("password2") {
					registerError = "The two passwords do not match"
					} else if _, err := services.GetUserID(r.FormValue("username")); err == nil {
						registerError = "The username is already taken"
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
		if err != nil {
			log.WithField("err", err).Error("Hashing error in PostRegisterUserHandler")
			return
		}
		username := r.FormValue("username")
		email := r.FormValue("email")
		user := models.User{Username: username, Email: email, PwHash: string(hash)}
		error := services.CreateUser(user)

		if error != nil {
			registerError = "Error while creating user"
		}

		AddFlash(session, w, r, "You are now registered ?", username)
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	data := models.PageData{
		"error": registerError,
	}
	redirectToRegister(w, data)
}

func redirectToRegister(w http.ResponseWriter, data models.PageData) {
	tmpl := LoadTemplate(RegisterPath)
	_ = tmpl.Execute(w, data)
}
