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

// PageData defines data on page whatever and request
type PageData map[string]interface{}

type Message struct {
	Email    string
	Username string
	Text     string
	PubDate  string
}

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

// query of the database just as above, but only finding us a single row
func queryDbSingleRow(query string, args ...interface{}) *sql.Row {
	liteDB, _ := sql.Open("sqlite3", database)

	res := liteDB.QueryRow(query, args...)

	return res
}

// getUserID returns user ID for username
func getUserID(username string) (int, error) {
	var ID int
	err := db.QueryRow("select user_id from user where username = ?", username).Scan(&ID)
	return ID, err
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
	res := queryDbSingleRow("select * from user where user_id = ?", userID)
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
			session.Values["username"] = tmpUser.username
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
	session, _ := store.Get(r, "_cookie")
	if session.Values["user_id"] != nil {
		routeName := fmt.Sprintf("/%s", session.Values["username"])
		http.Redirect(w, r, routeName, http.StatusFound)
	}

	http.Redirect(w, r, "/public", http.StatusFound)
}

func publicTimelineHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(timelinePath, layoutPath)
	if err != nil {
		log.Fatal(err)
	}
	res := queryDb("select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit ?", perPage)
	var msgs []Message

	for res.Next() {
		var (
			messageID int
			authorID  int
			text      string
			pubDate   int
			flagged   int
			userID    int
			username  string
			email     string
			pwHash    string
		)

		err = res.Scan(&messageID, &authorID, &text, &pubDate, &flagged, &userID, &username, &email, &pwHash)

		if err != nil {
			log.Fatal(err)
		}
		msgs = append(msgs, Message{
			Text:     text,
			PubDate:  formatDatetime(int64(pubDate)),
			Username: username,
			Email:    gravatarURL(email, 48),
		})
	}

	session, _ := store.Get(r, "_cookie")

	data := PageData{
		"username": session.Values["username"],
		"messages": msgs,
		"msgCount": len(msgs),
	}

	tmpl.Execute(w, data)
}

func userTimelineHandler(w http.ResponseWriter, r *http.Request) {
	//Display's a users tweets.
	params := mux.Vars(r)
	profileUsername := params["username"]

	profileUserID, err := getUserID(profileUsername)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	followed := false

	session, err := store.Get(r, "_cookie")
	if ok := session.Values["user_id"] != nil; ok {
		sessionUserID := session.Values["user_id"]      //Retrieves their username
		res := queryDb("select 1 from follower where "+ //Determines if the signed in user is following the user being viewed?
			"follower.who_id = ? and follower.whom_id = ?",
			sessionUserID, profileUserID)
		followed = res.Next() // Checks if the user that is signed in, is currently following the user on the page
	}

	tmpl, err := template.ParseFiles(timelinePath, layoutPath)
	if err != nil {
		log.Fatal(err)
	}
	messagesAndUsers, err := db.Query("select message.*, user.* from message, user where "+
		"user.user_id = message.author_id and user.user_id = ? "+
		"order by message.pub_date desc limit ?",
		profileUserID, perPage)

	data := PageData{"followed": followed}
	var messages []Message
	if err == nil {
		for messagesAndUsers.Next() {
			var (
				messageID int
				authorID  int
				text      string
				pubDate   int
				flagged   int
				userID    int
				username  string
				email     string
				pwHash    string
			)
			err := messagesAndUsers.Scan(&messageID, &authorID, &text, &pubDate, &flagged, &userID, &username, &email, &pwHash)
			if err == nil {
				message := Message{
					Email:    gravatarURL(email, 48),
					Username: username,
					Text:     text,
					PubDate:  formatDatetime(int64(pubDate)),
				}
				messages = append(messages, message)
			}
		}
	}
	data["messages"] = messages
	data["title"] = fmt.Sprintf("%s's Timeline", profileUsername)
	data["profileOwner"] = profileUsername
	data["followed"] = false

	if session.Values["username"] == profileUsername {
		data["ownProfile"] = true
	} else {
		currentUser := session.Values["user_id"]
		otherUser, err := getUserID(profileUsername)
		if err != nil {
			http.Error(w, "User does not exist", 400)
			return
		}
		res := queryDbSingleRow("select 1 from follower where who_id= ? and whom_id= ?", otherUser, currentUser)
		var (
			whoID  int
			whomID int
		)
		res.Scan(&whoID, &whomID)

		if whoID != 0 && whomID != 0 {
			data["followed"] = true
		}
	}

	data["msgCount"] = len(messages)
	data["username"] = session.Values["username"]

	tmpl.Execute(w, data)
}

func followUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "_cookie")
	currentUserID := session.Values["user_id"]
	params := mux.Vars(r)
	username := params["username"]
	userToFollowID, _ := getUserID(username)

	statement, err := db.Prepare(`insert into follower (who_id,whom_id) values(?,?)`)
	if err != nil {
		log.Fatal(err)
		return
	}
	statement.Exec(currentUserID, userToFollowID)
	statement.Close()
	routeName := fmt.Sprintf("/%s", username)
	http.Redirect(w, r, routeName, http.StatusFound)
}

// relies on a query string

func unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "_cookie")
	loggedInUser := session.Values["user_id"]
	var unfollowError string // keeping this so as to able to display an error message on the timeline
	// if we wanted one there
	if session.Values["user_id"] == nil {
		unfollowError = "no user was logged in "
		fmt.Println(unfollowError)
		http.Redirect(w, r, "timeline", http.StatusFound)
		return
	}
	v := r.URL.Query()
	p := v.Get("user")
	if p == "" {
		unfollowError = "the query parameter is empty"
		fmt.Println(unfollowError)
		http.Redirect(w, r, "timeline", http.StatusFound)
		return
	}
	if p == "" {
		unfollowError = "the query parameter is empty"
		fmt.Println(unfollowError)
		http.Redirect(w, r, "timeline", http.StatusFound)
		return
	}

	id2, user2Err := getUserID(p)

	if user2Err != nil {
		unfollowError = "no such user "
		fmt.Println(unfollowError)
		http.Redirect(w, r, "timeline", http.StatusFound)
		return
	}

	statement, er := db.Prepare("delete from follower where who_id= ? and whom_id= ?") // Prepare statement.

	if er != nil {
		fmt.Println("fuck fuck")
	}

	_, error := statement.Exec(id2, loggedInUser)
	statement.Close()
	if error != nil {
		unfollowError = "error during database operation "
		fmt.Println(unfollowError)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

}

func addMessageHandler(w http.ResponseWriter, r *http.Request) {}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "_cookie")
	if ok := session.Values["user_id"] != nil; ok {
		routeName := fmt.Sprintf("/%s", session.Values["username"])
		http.Redirect(w, r, routeName, http.StatusFound)
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

			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	tmpl, err := template.ParseFiles(loginPath, layoutPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	data := PageData{
		"error":    loginError,
		"username": session.Values["username"],
	}
	tmpl.Execute(w, data)

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "_cookie")
	if ok := session.Values["user_id"] != nil; ok {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	var registerError string
	if r.Method == "POST" {
		if len(r.FormValue("username")) == 0 {
			registerError = "You have to enter a username"
		} else if len(r.FormValue("email")) == 0 || strings.Contains(r.FormValue("email"), "@") == false {
			registerError = "You have to enter a valid email address"
		} else if len(r.FormValue("password")) == 0 {
			registerError = "You have to enter a password"
		} else if r.FormValue("password") != r.FormValue("password2") {
			registerError = "The two passwords do not match"
		} else if _, err := getUserID(r.FormValue("username")); err == nil {
			registerError = "The username is already taken"
		} else {
			statement, err := db.Prepare(`insert into user (username, email, pw_hash) values(?,?,?)`)
			if err != nil {
				log.Fatal(err)
				return
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal(err)
				return
			}

			statement.Exec(r.FormValue("username"), r.FormValue("email"), hash)
			statement.Close()
			session.AddFlash("You are now registered ?", r.FormValue("username"))
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}

	t, err := template.ParseFiles(registerPath, layoutPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		"error": registerError,
	}
	t.Execute(w, data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "_cookie")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	session.Values["user_id"] = ""
	session.Values["username"] = ""
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// init is automatically executed on program startup. Can't be called
// or referenced.
func init() {
	database, err := connectDb()
	if err != nil {
		log.Fatal(err)
	}
	db = database
}

func main() {
	router := mux.NewRouter()

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(s)

	router.Use(beforeRequest)
	router.Use(afterRequest)
	router.HandleFunc("/", timelineHandler)
	router.HandleFunc("/{username}/follow", followUserHandler)
	router.HandleFunc("/{username}/unfollow", unfollowUserHandler)
	router.HandleFunc("/login", loginHandler).Methods("GET", "POST")
	router.HandleFunc("/logout", logoutHandler)
	router.HandleFunc("/register", registerHandler).Methods("GET", "POST")
	router.HandleFunc("/public", publicTimelineHandler)
	router.HandleFunc("/{username}", userTimelineHandler)

	fmt.Println("Server running on port http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
