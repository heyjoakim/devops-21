package ui

import "net/http"

// FaviconHandler serves the site's favicon
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/dev.png")
}
