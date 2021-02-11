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
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// PageData defines data on page whatever
type PageData map[string]interface{}

type User struct {
	userID   int
	username string
	email    string
	pwHash   string
}

type layoutPage struct {
	Layout string
}

// configuration
var (
	database  = "./tmp/minitwit.db"
	perPage   = 30
	debug     = true
	secretKey = []byte("development key")
	store     = sessions.NewCookieStore(secretKey)
)

var db *sql.DB
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

func getUser(userID int) User {
	var (
		ID       int
		username string
		email    string
		pwHash   string
	)
	res := queryDb("select * from user where user_id = ?", userID)
	res.Scan(&ID, &username, &email, &pwHash)

	return User{
		userID:   ID,
		username: username,
		email:    email,
		pwHash:   pwHash,
	}
}

// beforeRequest make sure we are connected to the database each request and look
// up the current user so that we know he's there.
func beforeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		session, _ := store.Get(r, "_cookie")
		userID := session.Values["user_id"]
		if userID != nil {
			tmpUser := getUser(userID.(int))
			session.Values["user_id"] = tmpUser.userID
			session.Save(r, w)
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Closes the database again at the end of the request.
func afterRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		fmt.Println("Entered: " + r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
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
	session, err := store.Get(r, "_cookie")
	if ok := session.Values["user_id"] != nil; ok {
		http.Redirect(w, r, "timeline", http.StatusFound)
	}

	var loginError string
	if r.Method == "POST" {
		result := queryDb("select * from user where username = ?", r.FormValue("username"), true)
		if !result.Next() {
			loginError = "Invalid username"
			return
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
			session.AddFlash("You were logged in")
			session.Values["user_id"] = userID
			session.Save(r, w)

			http.Redirect(w, r, "/timeline", http.StatusFound)
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
	db, err := connectDb()
	name := r.FormValue("username")
	fmt.Println(name)
	fmt.Println(db)

	log.Println(r.Method)

	fmt.Println("Error")

	if r.Method == "POST" {
		// Error handling
		fmt.Println("HEP")
		test, _ := db.Query("select * from user where username = svopper1", true)
		fmt.Println(test)
		statement, _ := db.Prepare(`insert into user (username, email, pw_hash) values(?,?,?)`)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("HEP2")
		statement.Exec(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
		statement.Close()

	} else {
		fmt.Println("Error method not POST")

	}

	// Only implemented for display NO LOGIC (Joakim)
	t, err := template.ParseFiles(registerPath, layoutPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		"title": "Minitwit",
	}
	t.Execute(w, data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {}

// init is automatically executed on program startup. Can't be called
// or referenced.
func init() {
	database, err := connectDb()
	if err != nil {
		panic(err)
	}
	db = database
}

func main() {
	router := mux.NewRouter()

	//router.HandleFunc("/", layoutHandler)

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(s)

	//router.Use(beforeRequest)
	//router.Use(afterRequest)

	router.HandleFunc("/", timelineHandler)
	router.HandleFunc("/{username}/follow", followUserHandler)
	router.HandleFunc("/login", loginHandler).Methods("GET", "POST")
	router.HandleFunc("/register", registerHandler).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
