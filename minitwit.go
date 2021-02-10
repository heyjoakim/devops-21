package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// PageData defines data on page whatever
type PageData map[string]interface{}

type layoutPage struct {
	Layout string
}

// configuration
var (
	database  = "./tmp/minitwit.db"
	perPage   = 30
	debug     = true
	secretKey = "development key"
)

var staticPath string = "/static"
var cssPath string = "/css"
var timelinePath string = "./templates/timeline.html"
var layoutPath string = "./templates/layout.html"
var loginPath string = "./templates/login.html"
var registerPath string = "./templates/register.html"

// connectDb returns a new connection to the database.
func connectDb() (*sql.DB, error) {
	return sql.Open("sqlite3", database)
}

// initDb creates the database tables.
func initDb() {}

// queryDb queries the database and returns a list of dictionaries.
func queryDb(query string, args ...interface{}) *sql.Rows {
	liteDB, _ := sql.Open("sqlite3", database)

	res, _ := liteDB.Query(query, args...)

	return res
}

// getUserID returns user ID for username
func getUserID(username string) []interface{} {
	return nil
}

// formatDatetime formats a timestamp for display.
func formatDatetime(timestamp int64) string {
	timeObject := time.Unix(timestamp, 0)
	return timeObject.Format("2006-02-01 @ 02:04")
}

// gravatarURL return the gravatar image for the given email address.
func gravatarURL(email string, size int) string {
	encodedEmail := hex.EncodeToString([]byte(strings.ToLower(strings.TrimSpace(email))))
	hashedEmail := fmt.Sprintf("%x", sha256.Sum256([]byte(encodedEmail)))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon&s=%d", hashedEmail, size)
}

// beforeRequest make sure we are connected to the database each request and look
// up the current user so that we know he's there.
func beforeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		db, err := connectDb()
		if err != nil {
			panic(err)
		}
		fmt.Println(db)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Closes the database again at the end of the request.
func afterRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)

		log.Println("Done!")
	})
}

// timelineHandler a users timeline or if no user is logged in it will
// redirect to the public timeline.  This timeline shows the user's
// messages as well as all the messages of followed users.
func timelineHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(timelinePath, layoutPath)
	if err != nil {
		log.Fatal(err)
	}
	data := PageData{
		"title": "Minitwit",
	}
	tmpl.Execute(w, data)
	return
}

func publicTimelineHandler(w http.ResponseWriter, r *http.Request) {

}

func userTimelineHandler(w http.ResponseWriter, r *http.Request) {}

func followUserHandler(w http.ResponseWriter, r *http.Request) {
	// example on extract url params
	params := mux.Vars(r)
	username := params["username"]
	fmt.Println(username)
}

func unfollowUserHandler(w http.ResponseWriter, r *http.Request) {}

func addMessageHandler(w http.ResponseWriter, r *http.Request) {}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("_cookie")
	if err == nil {
		http.Redirect(w, r, "timeline", http.StatusFound)
	}

	var loginError string
	if r.Method == "POST" {
		result := queryDb("select * from user where username = ?", r.FormValue("username"), true)
		if !result.Next() {
			loginError = "Invalid username"

		}
		var (
			userID   int
			username string
			email    string
			pwHash   string
		)
		if err := result.Scan(&userID, &username, &email, &pwHash); err != nil {
			log.Fatal(err)
			loginError = "Invalid username"
		} else if err := bcrypt.CompareHashAndPassword([]byte(pwHash), []byte(r.FormValue("password"))); err != nil {
			loginError = "Invalid password"
		} else {
			cookie := http.Cookie{
				Name:     "_cookie",
				Value:    username,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		}
	}
	tmpl, err := template.ParseFiles(loginPath, layoutPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	data := PageData{
		"error": loginError,
	}
	tmpl.Execute(w, data)

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Only implemented for display NO LOGIC (Joakim)
	t, err := template.ParseFiles(registerPath, layoutPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	data := PageData{
		"title": "Minitwit",
	}
	t.Execute(w, data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	router := mux.NewRouter()

	//router.HandleFunc("/", layoutHandler)

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(s)

	router.Use(beforeRequest)
	router.Use(afterRequest)

	router.HandleFunc("/", timelineHandler)
	router.HandleFunc("/{username}/follow", followUserHandler)
	router.HandleFunc("/login", loginHandler).Methods("GET", "POST")
	router.HandleFunc("/register", registerHandler)

	log.Fatal(http.ListenAndServe(":8000", router))
}
