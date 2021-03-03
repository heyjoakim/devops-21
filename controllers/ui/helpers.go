package ui

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

// AddFlash add a flash to the session
func AddFlash(session *sessions.Session, w http.ResponseWriter, r *http.Request, message interface{}, vars ...string) {
	session.AddFlash(message, vars...)
	session.Save(r, w)
}

// LoadTemplate returns a HTML template
func LoadTemplate(path string) *template.Template {
	tmpl, err := template.ParseFiles(path, LayoutPath)
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}
