package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt time.Time
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
	prefillCaches() // _dev
	httpMultiplexer := http.NewServeMux()
	httpMultiplexer.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("www/styles"))))
	httpMultiplexer.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("www/js"))))
	httpMultiplexer.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("www/assets"))))

	httpMultiplexer.HandleFunc("/", serveWebpage)
	httpMultiplexer.HandleFunc("POST /login", loginUser)          // hx-action
	httpMultiplexer.HandleFunc("PUT /login", updateUser)          // hx-action
	httpMultiplexer.HandleFunc("POST /message", userPostsMessage) // hx-action
	httpMultiplexer.HandleFunc("GET /messages", refreshMessages)  // hx-action

	fmt.Println("Listening on 'http://localhost:8080/':")
	http.ListenAndServe(":8080", httpMultiplexer)
}

func serveWebpage(writer http.ResponseWriter, req *http.Request) {
	var tmpl *template.Template
	var data any

	/// active sessionToken? get user:
	userSession, _ := authenticateRequest(req)

	/// routing ///
	switch req.URL.Path {
	case "/flextest":
		tmpl = template.Must(template.ParseFiles("./www/flextest.html"))

	case "/myspace":
		tmpl = template.Must(template.ParseFiles(
			"./www/space.html",
			"./www/legos/nav.html",
		))
		cacheMutex.RLock()
		data = userCache[userSession.userId]
		cacheMutex.RUnlock()

	case "/":
		tmpl = template.Must(template.ParseFiles(
			"./www/index.html",
			"./www/legos/nav.html",
			"./www/legos/messageboard.html",
			"./www/legos/loginbar.html",
		))
		cacheMutex.RLock()
		data = struct {
			User     User
			Messages map[int]Message
		}{
			User:     userCache[userSession.userId],
			Messages: messageCache,
		}
		cacheMutex.RUnlock()

	default:
		http.Redirect(writer, req, "/", http.StatusSeeOther)
		return
	}
	fmt.Printf("- serving '%v' to client: %v\n", tmpl.Name(), req.RemoteAddr)
	tmpl.Execute(writer, data)
}
