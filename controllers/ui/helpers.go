package ui

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// AddFlash add a flash to the session
func AddFlash(session *sessions.Session, w http.ResponseWriter, r *http.Request, message interface{}, vars ...string) {
	session.AddFlash(message, vars...)
	session.Save(r, w)
}
