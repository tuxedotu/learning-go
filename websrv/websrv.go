package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type User struct { //type: 'User'
	Name string `json:"name"`
	Pass string `json:"pass"`
}

var userCache = make(map[int]User)
var userCacheMutex sync.RWMutex

func main() {
	httpMultiplexer := http.NewServeMux()
	httpMultiplexer.HandleFunc("/", handleRoot)           // web
	httpMultiplexer.HandleFunc("POST /users", createUser) // api

	fmt.Println("Listening on 'http://localhost:8080/':")
	http.ListenAndServe(":8080", httpMultiplexer)
}

func handleRoot(writer http.ResponseWriter, req *http.Request) {
	fmt.Printf("- serving '/' to %v\n", req.RemoteAddr)
	fmt.Fprintln(writer, "Hello Browser")
}

func createUser(writer http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(writer, "Name required!", http.StatusBadRequest)
		return
	}

	userCacheMutex.Lock()              // locking all other threads for time of userCache edits to prevent race conditions
	userCache[len(userCache)+1] = user // add new 'user' at endpos+1
	userCacheMutex.Unlock()

	writer.WriteHeader(http.StatusOK)
}
