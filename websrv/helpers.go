package main

import (
	"time"
)

// TODO: handleSessionToken(s) -> get, set, update, del(?)

func prefillCaches() { //>> test data
	userCache[1] = User{
		Id:        1,
		Name:      "admin",
		CreatedAt: time.Now(),
	}

	messageCache[1] = Message{
		Id:        1,
		Author:    "admin",
		Msg:       "hello",
		CreatedAt: time.Now(),
	}
	messageCache[2] = Message{
		Id:        2,
		Author:    "admin",
		Msg:       "messages!",
		CreatedAt: time.Now(),
	}
}
