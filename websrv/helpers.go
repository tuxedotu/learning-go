package main

import (
	"errors"
	"net/http"
	"time"
)

// TODO: handleSessionToken(s) -> get, set, update, del(?)

func authenticateRequest(req *http.Request) (session, error) {
	var retSesh session
	var err error

	if len(req.CookiesNamed("session_token")) < 1 {
		err = errors.New("Unauthenticated!")
		return retSesh, err
	}
	sessionToken := req.CookiesNamed("session_token")[0].Value

	cacheMutex.Lock()
	if time.Now().After(sessionsCache[sessionToken].expiry) {
		err = errors.New("Invalid session-token!")
		cacheMutex.Unlock()
		return retSesh, err
	}

	retSesh = sessionsCache[sessionToken]
	cacheMutex.Unlock()
	// fmt.Printf("- (authReq) LOG: > session = '%v' \n > error = '%v'", retSesh.expiry, err)
	return retSesh, err
}

func prefillCaches() { //>> test data
	userCache[1] = User{
		Id:        1,
		Name:      "admin",
		CreatedAt: time.Now(),
	}

	messageCache[1] = Message{
		Id:        1,
		Author:    "admin",
		Msg:       "hello, messages!",
		CreatedAt: time.Now(),
	}
}
