package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

func createUser(writer http.ResponseWriter, req *http.Request) {
	var user User
	var sessionToken string

	/// add new user object based on form input ///
	user = User{Name: req.FormValue("name"), CreatedAt: time.Now()} // TODO: potentially unsafe? (raw user form-input)
	if user.Name == "" {
		fmt.Println("- ERR in 'createUser': Name required!")
		http.Error(writer, "Name required!", http.StatusBadRequest)
		return
	}
	cacheMutex.Lock() //>> locking all other threads for time of edits to userCache to prevent race conditions
	user.Id = len(userCache) + 1
	userCache[user.Id] = user //>> add new 'user' at endpos+1
	cacheMutex.Unlock()
	fmt.Printf("- POST user (OK): %v\n", user)

	/// based on new user -> create new sessionToken & add cookie ///
	cmdOut, err := exec.Command("uuidgen").Output()
	if err != nil {
		fmt.Printf("- ERR: %v\n", err)
		return
	}
	sessionToken = string(cmdOut)[:len(string(cmdOut))-1] //>> byte-array to clean string conversion (w/o '\n')
	sessionsCache[sessionToken] = session{
		userId:    user.Id,
		expiry:    user.CreatedAt.Add(1 * time.Minute),
		createdAt: time.Now(),
	}
	http.SetCookie(writer, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: sessionsCache[sessionToken].expiry,
	})
	fmt.Printf("- session created: %v >> %T\n", sessionToken, sessionsCache[sessionToken])

	/// write response ///
	writer.Header().Set("HX-Redirect", "/") //>> *magic* HTMX-Page Redirects
	writer.WriteHeader(http.StatusOK)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var userInput string
	var familiarUser User

	userInput = r.FormValue("name") //> unsanitary user-input ??
	for _, user := range userCache {
		if user.Name == userInput {
			familiarUser = user
		}
	}
	fmt.Printf("- LOG: login attempt; input = '%v' (u: %v)\n", userInput, familiarUser.Name)

	if familiarUser == (User{}) {
		fmt.Printf("- WARN: failed login attempt; no such user!\n")
		http.Error(w, "No such user!", http.StatusBadRequest)
		w.Write([]byte("No such user!"))
		return
	}

	for sessionToken, session := range sessionsCache {
		if session.userId == familiarUser.Id && time.Now().After(session.expiry) {
			session.expiry = time.Now().Add(1 * time.Minute)
			sessionsCache[sessionToken] = session
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: session.expiry,
			})
		}
	}
	if w.Header().Get("session_token") == "" {
		cmdOut, err := exec.Command("uuidgen").Output()
		if err != nil {
			fmt.Printf("- ERR: couldn't create sessionToken!")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		sessionToken := string(cmdOut)[:len(string(cmdOut))-1] //>> byte-array to clean string conversion (w/o '\n')
		newSession := session{
			userId:    familiarUser.Id,
			expiry:    time.Now().Add(1 * time.Minute),
			createdAt: time.Now(),
		}
		sessionsCache[string(sessionToken)] = newSession
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   string(sessionToken),
			Expires: newSession.expiry,
		})
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var userInput string
	var session session
	var user User

	userInput = r.FormValue("name")
	if userInput == "" {
		fmt.Printf("- ERR: no empty user-updates allowed!\n")
		http.Error(w, "No empty updates allowed!", http.StatusBadRequest)
		return
	}
	if len(r.CookiesNamed("session_token")) < 1 {
		fmt.Printf("- ERR: trying update w/o valid token!\n")
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session, ok := sessionsCache[r.CookiesNamed("session_token")[0].Value]
	if !ok {
		fmt.Printf("- ERR: no such token '%v'\n", r.CookiesNamed("session_token")[0].Value)
		http.Error(w, "Invalid session token!", http.StatusBadRequest)
		return
	}
	cacheMutex.Lock()
	user = userCache[session.userId]
	user.Name = userInput
	userCache[session.userId] = user
	cacheMutex.Unlock()

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func userPostsMessage(w http.ResponseWriter, r *http.Request) {
	var session session
	var userInput string
	var postedMessage Message

	userInput = r.FormValue("message")
	if userInput == "" {
		fmt.Printf("- ERR: no empty user-updates allowed!\n")
		http.Error(w, "No empty updates allowed!", http.StatusBadRequest)
		return
	}
	if len(r.CookiesNamed("session_token")) < 1 {
		fmt.Printf("- ERR: trying update w/o valid token!\n")
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session, ok := sessionsCache[r.CookiesNamed("session_token")[0].Value]
	if !ok {
		fmt.Printf("- ERR: no such token '%v'\n", r.CookiesNamed("session_token")[0].Value)
		http.Error(w, "Invalid session token!", http.StatusBadRequest)
		return
	}
	postedMessage = Message{
		Id:        len(messageCache) + 1,
		Author:    userCache[session.userId].Name,
		Msg:       userInput,
		CreatedAt: time.Now(),
	}
	fmt.Printf("- LOG: postedMessage = '%v'\n", postedMessage)
	cacheMutex.Lock()
	messageCache[postedMessage.Id] = postedMessage
	cacheMutex.Unlock()

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func refreshMessages(w http.ResponseWriter, r *http.Request) {
	var messageList = []struct {
		Author string
		Msg    string
	}{}
	var tmpl *template.Template

	fmt.Printf("- LOG: Msg-list = '%v'\n", messageList)
	tmpl = template.Must(tmpl.ParseFiles("./www/legos/messageboard.html"))
	tmpl.Execute(w, messageList)
}

// NOT USED: just an api-style json example //
func getUser(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	cacheMutex.RLock()
	user, ok := userCache[id]
	cacheMutex.RUnlock()
	if !ok {
		http.Error(writer, "Bad Request!", http.StatusBadRequest)
		return
	}

	jsonUser, err := json.Marshal(user)
	if !ok {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonUser)
}
