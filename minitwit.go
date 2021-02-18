package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/heyjoakim/devops-21/models"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// PageData defines data on page whatever and request
type PageData map[string]interface{}

// Message defines message
type Message struct {
	Email    string
	Username string
	Text     string
	PubDate  string
}

// User defines a user
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

// var db *sql.DB

// App defines the application
type App struct {
	db *sql.DB
}

var staticPath string = "/static"
var cssPath string = "/css"
var timelinePath string = "./templates/timeline.html"
var layoutPath string = "./templates/layout.html"
var loginPath string = "./templates/login.html"
var registerPath string = "./templates/register.html"

// connectDb returns a new connection to the database.
func (d *App) connectDb() (*sql.DB, error) {
	return sql.Open("sqlite3", database)
}

// initDb creates the database tables.
func (d *App) initDb() {
	file, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		log.Print(err.Error())
	}
	tx, _ := d.db.Begin()
	_, err = d.db.Exec(string(file))
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("Unable to rollback initDb: %v", rollbackErr)
		}
		log.Fatal(err)
	}
}

// queryDb queries the database and returns a list of dictionaries.
func (d *App) queryDb(query string, args ...interface{}) *sql.Rows {
	res, _ := d.db.Query(query, args...)
	return res
}

// query of the database just as above, but only finding us a single row
func (d *App) queryDbSingleRow(query string, args ...interface{}) *sql.Row {
	liteDB, _ := sql.Open("sqlite3", database)

	res := liteDB.QueryRow(query, args...)

	return res
}

// getUserID returns user ID for username
func (d *App) getUserID(username string) (int, error) {
	var ID int
	err := d.db.QueryRow("select user_id from user where username = ?", username).Scan(&ID)
	return ID, err
}

// formatDatetime formats a timestamp for display.
func (d *App) formatDatetime(timestamp int64) string {
	timeObject := time.Unix(timestamp, 0)
	return timeObject.Format("2006-02-01 @ 02:04")
}

// gravatarURL return the gravatar image for the given email address.
func (d *App) gravatarURL(email string, size int) string {
	encodedEmail := hex.EncodeToString([]byte(strings.ToLower(strings.TrimSpace(email))))
	hashedEmail := fmt.Sprintf("%x", sha256.Sum256([]byte(encodedEmail)))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon&s=%d", hashedEmail, size)
}

func (d *App) getUser(userID int) models.User {
	var (
		ID       int
		username string
		email    string
		pwHash   string
	)
	res := d.queryDbSingleRow("select * from user where user_id = ?", userID)
	res.Scan(&ID, &username, &email, &pwHash)

	return models.User{
		UserID:   ID,
		Username: username,
		Email:    email,
		PwHash:   pwHash,
	}
}

// beforeRequest make sure we are connected to the database each request and look
// up the current user so that we know he's there.
func (d *App) beforeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		database, _ := d.connectDb()
		d.db = database
		session, _ := store.Get(r, "_cookie")
		userID := session.Values["user_id"]
		if userID != nil {
			tmpUser := d.getUser(userID.(int))
			session.Values["user_id"] = tmpUser.UserID
			session.Values["username"] = tmpUser.Username
			session.Save(r, w)
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Closes the database again at the end of the request.
func (d *App) afterRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Entered: " + r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		d.db.Close()
	})
}

// timelineHandler a users timeline or if no user is logged in it will
// redirect to the public timeline.  This timeline shows the user's
// messages as well as all the messages of followed users.
func (d *App) timelineHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "_cookie")

	if session.Values["user_id"] != nil {
		routeName := fmt.Sprintf("/%s", session.Values["username"])
		http.Redirect(w, r, routeName, http.StatusFound)
		return
	}

	http.Redirect(w, r, "/public", http.StatusFound)
}

func (d *App) publicTimelineHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(timelinePath, layoutPath)
	if err != nil {
		log.Fatal(err)
	}
	res := d.queryDb("select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit ?", perPage)
	var msgs []models.Message

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
		msgs = append(msgs, models.Message{
			Text:     text,
			PubDate:  d.formatDatetime(int64(pubDate)),
			Username: username,
			Email:    d.gravatarURL(email, 48),
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

func (d *App) userTimelineHandler(w http.ResponseWriter, r *http.Request) {
	//Display's a users tweets.
	params := mux.Vars(r)
	profileUsername := params["username"]

	profileUserID, err := d.getUserID(profileUsername)
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err)
		return
	}
	followed := false

	session, err := store.Get(r, "_cookie")
	if ok := session.Values["user_id"] != nil; ok {
		sessionUserID := session.Values["user_id"]        //Retrieves their username
		res := d.queryDb("select 1 from follower where "+ //Determines if the signed in user is following the user being viewed?
			"follower.who_id = ? and follower.whom_id = ?",
			sessionUserID, profileUserID)
		defer res.Close()
		followed = res.Next() // Checks if the user that is signed in, is currently following the user on the page
	}

	tmpl, err := template.ParseFiles(timelinePath, layoutPath)
	if err != nil {
		log.Fatal(err)
	}
	messagesAndUsers, err := d.db.Query("select message.*, user.* from message, user where "+
		"user.user_id = message.author_id and user.user_id = ? "+
		"order by message.pub_date desc limit ?",
		profileUserID, perPage)
	if err != nil {
		fmt.Println("Err retrieving messages", err)
	}

	var msgS []models.Message
	if ok := session.Values["user_id"] != nil; ok {
		sessionUserID := session.Values["user_id"].(int)
		if sessionUserID == profileUserID {
			followlist := d.getFollowedUsers(sessionUserID)
			for _, v := range followlist {
				msgS = append(d.getPostsForuser(v), msgS...)
			}
		}
	}

	data := PageData{"followed": followed}
	var messages []models.Message
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
				message := models.Message{
					Email:    d.gravatarURL(email, 48),
					Username: username,
					Text:     text,
					PubDate:  d.formatDatetime(int64(pubDate)),
				}
				messages = append(messages, message)
			}
		}
	}
	messages = append(msgS, messages...)
	data["messages"] = messages
	data["title"] = fmt.Sprintf("%s's Timeline", profileUsername)
	data["profileOwner"] = profileUsername

	if session.Values["username"] == profileUsername {
		data["ownProfile"] = true
	} else {
		currentUser := session.Values["user_id"]
		otherUser, err := d.getUserID(profileUsername)
		if err != nil {
			http.Error(w, "User does not exist", 400)
			return
		}

		res := d.queryDbSingleRow("select 1 from follower where who_id= ? and whom_id= ?", otherUser, currentUser)
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
	data["MsgInfo"] = session.Flashes("Info")
	data["MsgWarn"] = session.Flashes("Warn")
	session.Save(r, w)

	tmpl.Execute(w, data)
}

func (d *App) getPostsForuser(id int) []models.Message {
	messagesAndUsers, err := d.db.Query("select message.*, user.* from message, user where "+
		"user.user_id = message.author_id and user.user_id = ? "+
		"order by message.pub_date desc limit ?",
		id, perPage)
	var messages []models.Message
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
				message := models.Message{
					Email:    d.gravatarURL(email, 48),
					Username: username,
					Text:     text,
					PubDate:  d.formatDatetime(int64(pubDate)),
				}
				messages = append(messages, message)
			}
		}
	}
	return messages
}

// get ID's of all users that are followed by some user
func (d *App) getFollowedUsers(id int) []int {
	followedIDs, _ := d.db.Query("select * from follower where  who_id= ?", id)
	var followlist []int
	for followedIDs.Next() {
		var (
			from int
			to   int
		)
		err := followedIDs.Scan(&from, &to)
		if err == nil {
			followID := to
			followlist = append(followlist, followID)
		}
	}

	return followlist
}
func (d *App) followUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "_cookie")
	currentUserID := session.Values["user_id"]
	params := mux.Vars(r)
	username := params["username"]
	userToFollowID, _ := d.getUserID(username)

	statement, err := d.db.Prepare(`insert into follower (who_id,whom_id) values(?,?)`)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, error := statement.Exec(currentUserID, userToFollowID)
	statement.Close()
	if error != nil {
		fmt.Println("database error: ", error)
	}

	routeName := fmt.Sprintf("/%s", username)
	session.AddFlash("You are now following "+username, "Info")
	session.Save(r, w)

	http.Redirect(w, r, routeName, http.StatusFound)
}

// relies on a query string

func (d *App) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "_cookie")
	loggedInUser := session.Values["user_id"]
	params := mux.Vars(r)
	username := params["username"]
	if username == "" {
		session.AddFlash("No query parameter present", "Warn")
		session.Save(r, w)
		http.Redirect(w, r, "timeline", http.StatusFound)
		return
	}

	id2, user2Err := d.getUserID(username)
	if user2Err != nil {
		session.AddFlash("User does not exist", "Warn")
		session.Save(r, w)
		http.Redirect(w, r, "timeline", http.StatusFound)
		return
	}

	statement, _ := d.db.Prepare("delete from follower where who_id= ? and whom_id= ?") // Prepare statement.
	_, error := statement.Exec(loggedInUser, id2)
	statement.Close()
	if error != nil {
		session.AddFlash("Error following user", "Warn")
		session.Save(r, w)
		fmt.Println("db error: ", error)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.AddFlash("You are no longer following "+username, "Info")
	session.Save(r, w)
	http.Redirect(w, r, "/"+username, http.StatusFound)
	return
}

func (d *App) addMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		statement, err := d.db.
			Prepare(`insert into message (author_id, text, pub_date, flagged) values(?,?,?,0)`)

		if err != nil {
			fmt.Println("db error during message creation") // probably needing some error handling
			log.Fatal(err)
			return
		}

		userID, _ := d.getUserID(r.FormValue("token"))
		statement.Exec(userID, r.FormValue("text"), time.Now().Unix())
		statement.Close()

		http.Redirect(w, r, "/"+r.FormValue("token"), http.StatusFound)
	}
}

func (d *App) loginHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "_cookie")
	if ok := session.Values["user_id"] != nil; ok {
		routeName := fmt.Sprintf("/%s", session.Values["username"])
		http.Redirect(w, r, routeName, http.StatusFound)
	}

	var (
		userID   int
		username string
		email    string
		pwHash   string
	)

	var loginError string
	if r.Method == "POST" {
		err := d.db.QueryRow("select * from user where username = ?", r.FormValue("username")).Scan(&userID, &username, &email, &pwHash)
		if err != nil {
			loginError = "User does not exist"
		}

		if r.FormValue("username") != username {
			loginError = "Invalid username"
		} else if err := bcrypt.CompareHashAndPassword([]byte(pwHash), []byte(r.FormValue("password"))); err != nil {
			loginError = "Invalid password"
		} else {
			session.AddFlash("You were logged in")
			session.Values["user_id"] = userID
			session.Save(r, w)

			http.Redirect(w, r, "/"+username, http.StatusFound)
		}
		// d.db.Close()
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

func (d *App) registerHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "_cookie")
	if ok := session.Values["user_id"] != nil; ok {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	var registerError string
	if r.Method == "POST" {
		if len(r.FormValue("username")) == 0 {
			registerError = "You have to enter a username"
		} else if len(r.FormValue("email")) == 0 || strings.
			Contains(r.FormValue("email"), "@") == false {
			registerError = "You have to enter a valid email address"
		} else if len(r.FormValue("password")) == 0 {
			registerError = "You have to enter a password"
		} else if r.FormValue("password") != r.FormValue("password2") {
			registerError = "The two passwords do not match"
		} else if _, err := d.getUserID(r.FormValue("username")); err == nil {
			registerError = "The username is already taken"
		} else {
			statement, err := d.db.
				Prepare(`insert into user (username, email, pw_hash) values(?,?,?)`)

			if err != nil {
				log.Fatal(err)
				return
			}

			hash, err := bcrypt.
				GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
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

func (d *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
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
func (d *App) init() {
	database, err := d.connectDb()
	if err != nil {
		log.Fatal(err)
	}
	d.db = database
}

func main() {
	router := mux.NewRouter()

	var app App

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(s)

	router.Use(app.beforeRequest)
	router.Use(app.afterRequest)
	router.HandleFunc("/", app.timelineHandler)
	router.HandleFunc("/{username}/unfollow", app.unfollowUserHandler)
	router.HandleFunc("/{username}/follow", app.followUserHandler)
	router.HandleFunc("/login", app.loginHandler).Methods("GET", "POST")
	router.HandleFunc("/logout", app.logoutHandler)
	router.HandleFunc("/addMessage", app.addMessageHandler).Methods("GET", "POST")
	router.HandleFunc("/register", app.registerHandler).Methods("GET", "POST")
	router.HandleFunc("/public", app.publicTimelineHandler)
	router.HandleFunc("/{username}", app.userTimelineHandler)

	fmt.Println("Server running on port http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
