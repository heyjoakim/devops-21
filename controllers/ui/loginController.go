package ui

import (
	"fmt"
	"net/http"

	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// GetLoginHandler returns the login page
func GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	session := GetSession(w, r)
	if ok := session.Values["user_id"] != nil; ok {
		routeName := fmt.Sprintf("/%s", session.Values["username"])
		http.Redirect(w, r, routeName, http.StatusFound)
	}

	data := models.PageData{
		"username": session.Values["username"],
	}
	redirectToLogin(w, data)
}

// PostLoginHandler handles user login
func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	session := GetSession(w, r)

	var loginError string
	user, err := services.GetUserFromUsername(r.FormValue("username"))
	if err != nil {
		loginError = "Unknown username"
		data := models.PageData{
			"error": loginError,
		}
		redirectToLogin(w, data)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(r.FormValue("password"))); err != nil {
		fmt.Println([]byte(user.PwHash))
		fmt.Println([]byte(r.FormValue("user")))

		loginError = "Invalid password"
		log.WithField("err", err).Error("Password hashes are different")
		data := models.PageData{
			"error":    loginError,
			"username": session.Values["username"],
		}
		redirectToLogin(w, data)
		return
	}
	session.Values["user_id"] = user.UserID
	AddFlash(session, w, r, "You were logged in")
	http.Redirect(w, r, "/"+user.Username, http.StatusFound)
}

func redirectToLogin(w http.ResponseWriter, data models.PageData) {
	tmpl := LoadTemplate(LoginPath)
	_ = tmpl.Execute(w, data)
}
