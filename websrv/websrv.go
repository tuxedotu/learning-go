package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"
)

type User struct { //type: 'User'
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	Id        int
	Author    string
	Msg       string
	CreatedAt time.Time
}

type session struct {
	userId    int
	expiry    time.Time
	createdAt time.Time
}

var userCache = make(map[int]User)
var sessionsCache = make(map[string]session)
var messageCache = make(map[int]Message)
var cacheMutex sync.RWMutex

func main() {
	prefillCaches()
	httpMultiplexer := http.NewServeMux()
	httpMultiplexer.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("www/styles"))))
	httpMultiplexer.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("www/js"))))
	httpMultiplexer.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("www/assets"))))

	httpMultiplexer.HandleFunc("/", serveWebpage)
	httpMultiplexer.HandleFunc("POST /login", loginUser)             // hx-action
	httpMultiplexer.HandleFunc("POST /signup", createUser)           // hx-action
	httpMultiplexer.HandleFunc("POST /update", updateUser)           // hx-action
	httpMultiplexer.HandleFunc("POST /message", userPostsMessage)    // hx-action
	httpMultiplexer.HandleFunc("GET /messageboard", refreshMessages) // hx-action

	fmt.Println("Listening on 'http://localhost:8080/':")
	http.ListenAndServe(":8080", httpMultiplexer)
}

func serveWebpage(writer http.ResponseWriter, req *http.Request) {
	var tmpl *template.Template
	var sessionToken string
	var familiarUser User
	var data any

	/// active sessionToken? get user:
	if len(req.CookiesNamed("session_token")) > 0 { //>> cookie is only attached before expiry, very nice!
		sessionToken = req.CookiesNamed("session_token")[0].Value

		if time.Now().Before(sessionsCache[sessionToken].expiry) {
			cacheMutex.RLock()
			familiarUser = userCache[sessionsCache[sessionToken].userId] //>> protected read of sessionToken-user
			cacheMutex.RUnlock()
		}
	} //>> TODO: delete old sessions from sessionsCache

	/// routing ///
	switch req.URL.Path {
	case "/login":
		tmpl = template.Must(template.ParseFiles(
			"./www/login.html",
			"./www/legos/nav.html",
		))
		data = familiarUser
	default:
		tmpl = template.Must(template.ParseFiles(
			"./www/index.html",
			"./www/legos/nav.html",
			"./www/legos/messageboard.html",
		))
		data = struct {
			User     User
			Messages map[int]Message
		}{
			User:     familiarUser,
			Messages: messageCache,
		}
	}

	fmt.Printf("- serving '%v' to client: %v\n", tmpl.Name(), req.RemoteAddr)
	tmpl.Execute(writer, data)
}
