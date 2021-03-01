package ui

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	// PerPage defines how many results are returned
	PerPage = 30
	store     = sessions.NewCookieStore(secretKey)
	secretKey = []byte("development key")
)

// GetSession returns the current browser session
func GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, err := store.Get(r, "_cookie")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	return session
}
