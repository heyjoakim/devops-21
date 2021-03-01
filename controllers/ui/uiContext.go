package ui

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	PerPage   = 30
	Debug     = true
	store     = sessions.NewCookieStore(secretKey)
	secretKey = []byte("development key")
)

func GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, err := store.Get(r, "_cookie")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	return session
}
