package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/heyjoakim/devops-21/models"
	services "github.com/heyjoakim/devops-21/services"
	"golang.org/x/crypto/bcrypt"
)

// configuration
var (
	perPage   = 30
	debug     = true
	secretKey = []byte("development key")
	store     = sessions.NewCookieStore(secretKey)
)

type layoutPage struct {
	Layout string
}

var (
	staticPath   string = "/static"
	cssPath      string = "/css"
	timelinePath string = "./templates/timeline.html"
	layoutPath   string = "./templates/layout.html"
	loginPath    string = "./templates/login.html"
	registerPath string = "./templates/register.html"
)

// // gravatarURL return the gravatar image for the given email address.
// func gravatarURL(email string, size int) string {
// 	encodedEmail := hex.EncodeToString([]byte(strings.ToLower(strings.TrimSpace(email))))
// 	hashedEmail := fmt.Sprintf("%x", sha256.Sum256([]byte(encodedEmail)))
// 	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon&s=%d", hashedEmail, size)
// }

// BeforeRequest make sure we are connected to the database each request and look
// up the current user so that we know he's there.
func BeforeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// database, _ := d.connectDb()
		// d.db = database
		session, _ := store.Get(r, "_cookie")
		userID := session.Values["user_id"]
		if userID != nil {
			id := userID.(uint)
			tmpUser := services.GetUser(id)
			session.Values["user_id"] = tmpUser.UserID
			session.Values["username"] = tmpUser.Username
			session.Save(r, w)
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Closes the database again at the end of the request.
func AfterRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("[%s] --> %s", r.Method, r.RequestURI))
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// // TimelineHandler a users timeline or if no user is logged in it will
// // redirect to the public timeline.  This timeline shows the user's
// // messages as well as all the messages of followed users.
// func TimelineHandler(w http.ResponseWriter, r *http.Request) {
// 	session, _ := store.Get(r, "_cookie")
// 	if session.Values["user_id"] != nil {
// 		routeName := fmt.Sprintf("/%s", session.Values["username"])
// 		http.Redirect(w, r, routeName, http.StatusFound)
// 		return
// 	}

// 	http.Redirect(w, r, "/public", http.StatusFound)
// }

// func PublicTimelineHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.ParseFiles(timelinePath, layoutPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var results = services.GetPublicMessages(perPage)

// 	var messages []models.MessageViewModel
// 	for _, result := range results {
// 		message := models.MessageViewModel{
// 			Email:   gravatarURL(result.Email, 48),
// 			User:    result.Username,
// 			Content: result.Content,
// 			PubDate: helpers.FormatDatetime(result.PubDate),
// 		}
// 		messages = append(messages, message)
// 	}

// 	session, err := store.Get(r, "_cookie")
// 	username := session.Values["username"]

// 	data := models.PageData{
// 		"username": username,
// 		"messages": messages,
// 		"msgCount": len(messages),
// 	}

// 	tmpl.Execute(w, data)
// }

// func UserTimelineHandler(w http.ResponseWriter, r *http.Request) {
// 	//Display's a users tweets.
// 	params := mux.Vars(r)
// 	profileUsername := params["username"]

// 	profileUserID, err := services.GetUserID(profileUsername)
// 	if err != nil {
// 		w.WriteHeader(404)
// 		fmt.Println(err)
// 		return
// 	}

// 	session, err := store.Get(r, "_cookie")
// 	sessionUserID := session.Values["user_id"]
// 	data := models.PageData{"followed": false}

// 	if sessionUserID != nil {
// 		if services.IsUserFollower(sessionUserID.(uint), profileUserID) {
// 			data["followed"] = true
// 		}
// 	}

// 	tmpl, err := template.ParseFiles(timelinePath, layoutPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	messages := getPostsForUser(profileUserID)

// 	var msgS []models.MessageViewModel
// 	if sessionUserID != nil {
// 		if sessionUserID == profileUserID {
// 			followlist := getFollowedUsers(sessionUserID.(uint))
// 			for _, v := range followlist {
// 				msgS = append(getPostsForUser(v), msgS...)
// 			}
// 		}
// 	}

// 	messages = append(msgS, messages...)
// 	data["messages"] = messages
// 	data["title"] = fmt.Sprintf("%s's Timeline", profileUsername)
// 	data["profileOwner"] = profileUsername

// 	if session.Values["username"] == profileUsername {
// 		data["ownProfile"] = true
// 	} /*else { //TODO Delete once it ahs been verified accessing a nonexistent userpage returns a 404
// 		currentUser := session.Values["user_id"]
// 		otherUser, err := services.GetUserID(profileUsername)
// 		if err != nil {
// 			http.Error(w, "User does not exist", 400)
// 			return
// 		}

// 		var follower models.Follower
// 		d.db.Where("who_id = ?", otherUser).
// 			Where("whom_id = ?", currentUser).
// 			First(&follower)
// 		if follower.WhoID != 0 && follower.WhomID != 0 {
// 			data["followed"] = true
// 		}
// 	}*/

// 	data["msgCount"] = len(messages)
// 	data["username"] = session.Values["username"]
// 	data["MsgInfo"] = session.Flashes("Info")
// 	data["MsgWarn"] = session.Flashes("Warn")
// 	session.Save(r, w)

// 	tmpl.Execute(w, data)
// }

// func getPostsForUser(id uint) []models.MessageViewModel {
// 	var results = services.GetMessagesForUser(perPage, id)

// 	var messages []models.MessageViewModel
// 	for _, result := range results {
// 		message := models.MessageViewModel{
// 			Email:   helpers.GetGravatarURL(result.Email, 48),
// 			User:    result.Username,
// 			Content: result.Content,
// 			PubDate: helpers.FormatDatetime(result.PubDate),
// 		}
// 		messages = append(messages, message)
// 	}

// 	return messages
// }

// // get ID's of all users that are followed by some user
// func getFollowedUsers(id uint) []uint {
// 	var followers = services.GetUsersFollowedBy(id)

// 	var followlist []uint
// 	for _, follower := range followers {
// 		followlist = append(followlist, follower.WhomID)
// 	}

// 	return followlist
// }

// // FollowUserHandler handles following another user
// func FollowUserHandler(w http.ResponseWriter, r *http.Request) {
// 	session, _ := store.Get(r, "_cookie")
// 	currentUserID := session.Values["user_id"].(uint)
// 	params := mux.Vars(r)
// 	username := params["username"]
// 	userToFollowID, _ := services.GetUserID(username)

// 	follower := models.Follower{WhoID: currentUserID, WhomID: userToFollowID}
// 	err := services.CreateFollower(follower)
// 	if err != nil {
// 		log.Print(err)
// 		fmt.Println("database error: ", err)
// 		return
// 	}

// 	session.AddFlash("You are now following "+username, "Info")
// 	session.Save(r, w)
// 	http.Redirect(w, r, "/"+username, http.StatusFound)
// }

// // UnfollowUserHandler - relies on a query string
// func UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
// 	session, _ := store.Get(r, "_cookie")
// 	loggedInUser := session.Values["user_id"].(uint)
// 	params := mux.Vars(r)
// 	username := params["username"]

// 	if username == "" {
// 		session.AddFlash("No query parameter present", "Warn")
// 		session.Save(r, w)
// 		http.Redirect(w, r, "timeline", http.StatusFound)
// 		return
// 	}

// 	id2, user2Err := services.GetUserID(username)
// 	if user2Err != nil {
// 		session.AddFlash("User does not exist", "Warn")
// 		session.Save(r, w)
// 		http.Redirect(w, r, "timeline", http.StatusFound)
// 		return
// 	}

// 	err := services.UnfollowUser(loggedInUser, id2)

// 	if err != nil {
// 		session.AddFlash("Error following user", "Warn")
// 		session.Save(r, w)
// 		fmt.Println("db error: ", err)
// 		http.Redirect(w, r, "/", http.StatusFound)
// 		return
// 	}

// 	session.AddFlash("You are no longer following "+username, "Info")
// 	session.Save(r, w)
// 	http.Redirect(w, r, "/"+username, http.StatusFound)
// 	return
// }

func AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID, _ := services.GetUserID(r.FormValue("token"))
		message := models.Message{AuthorID: userID, Text: r.FormValue("text"), PubDate: time.Now().Unix(), Flagged: 0}
		err := services.CreateMessage(message)
		if err != nil {
			fmt.Println("database error: ", err)
			log.Fatal(err)
			return
		}

		http.Redirect(w, r, "/"+r.FormValue("token"), http.StatusFound)
	}
}

func RegisterUserUiHandler(w http.ResponseWriter, r *http.Request) {
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
		} else if _, err := services.GetUserID(r.FormValue("username")); err == nil {
			registerError = "The username is already taken"
		} else {

			hash, err := bcrypt.
				GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal(err)
				return
			}
			username := r.FormValue("username")
			email := r.FormValue("email")
			user := models.User{Username: username, Email: email, PwHash: string(hash)}
			error := services.CreateUser(user)

			if error != nil {
				fmt.Println(error)
				registerError = "Error while creating user"
			}

			session.AddFlash("You are now registered ?", username)
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}

	t, err := template.ParseFiles(registerPath, layoutPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := models.PageData{
		"error": registerError,
	}

	t.Execute(w, data)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
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

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/dev.png")
}
