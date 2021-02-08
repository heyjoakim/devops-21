package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

// PageData defines data on page whatever
type PageData map[string]string

// configuration
var (
	database  = "./tmp/minitwit.db"
	perPage   = 30
	debug     = true
	secretKey = "development key"
)

var staticPath string = "/static"
var indexPath string = "./templates/timeline.html"

// connectDb returns a new connection to the database.
func connectDb() interface{} { // replace interface return type with whatever golang sqlite lib returns
	return nil
}

// initDb creates the database tables.
func initDb() {}

// queryDb queries the database and returns a list of dictionaries.
func queryDb(query string, args interface{}, one bool) []interface{} { // replace []interface return type with whatever golang sqlite lib returns
	return nil
}

// getUserID returns user ID for username
func getUserID(username string) []interface{} {
	return nil
}

// formatDatetime formats a timestamp for display.
func formatDatetime(timestamp int64) string { return "" }

// gravatarURL return the gravatar image for the given email address.
func gravatarURL(email string, size int) string { return "" }

// beforeRequest make sure we are connected to the database each request and look
// up the current user so that we know he's there.
func beforeRequest() {}

// Closes the database again at the end of the request.
func afterRequest(respone interface{}) interface{} {
	return nil
}

// timelineHandler a users timeline or if no user is logged in it will
// redirect to the public timeline.  This timeline shows the user's
// messages as well as all the messages of followed users.
func timelineHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(indexPath)
	if err != nil {
		log.Fatal(err)
	}
	data := PageData{
		"title": "Minitwit",
	}
	tmpl.Execute(w, data)
	return
}

func publicTimelineHandler(w http.ResponseWriter, r *http.Request) {}

func userTimelineHandler(w http.ResponseWriter, r *http.Request) {}

func followUserHandler(w http.ResponseWriter, r *http.Request) {
	// example on extract url params
	params := mux.Vars(r)
	username := params["username"]
	fmt.Println(username)
}

func unfollowUserHandler(w http.ResponseWriter, r *http.Request) {}

func addMessageHandler(w http.ResponseWriter, r *http.Request) {}

func loginHandler(w http.ResponseWriter, r *http.Request) {}

func registerHandler(w http.ResponseWriter, r *http.Request) {}

func logoutHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", timelineHandler)
	router.HandleFunc("/{username}/follow", followUserHandler)
	log.Fatal(http.ListenAndServe(":8000", router))
}
