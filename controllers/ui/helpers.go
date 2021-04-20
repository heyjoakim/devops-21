package ui

import (
	"net/http"
	"text/template"

	"github.com/heyjoakim/devops-21/models"
	"github.com/heyjoakim/devops-21/services"
	"github.com/gorilla/sessions"
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
		services.LogError(models.Log{
			Message: err.Error(),
		})
	}
	return tmpl
}
