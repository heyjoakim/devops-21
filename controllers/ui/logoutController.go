package ui

import "net/http"

// LogoutHandler handles user logout. It removed any information related to the user.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session := GetSession(w, r)

	session.Values["user_id"] = ""
	session.Values["username"] = ""
	session.Options.MaxAge = -1
	_ = session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}
