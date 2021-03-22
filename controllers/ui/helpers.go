package ui

import (
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

// AddFlash add a flash to the session
func AddFlash(session *sessions.Session, w http.ResponseWriter, r *http.Request, message interface{}, vars ...string) {
	session.AddFlash(message, vars...)
	_ = session.Save(r, w)
}

// LoadTemplate returns a HTML template
func LoadTemplate(path string) *template.Template {
	tmpl, err := template.ParseFiles(path, LayoutPath)
	if err != nil {
		log.Error(err)
	}
	return tmpl
}
