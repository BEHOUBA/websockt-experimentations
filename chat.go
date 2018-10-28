package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ChatRoom struct {
	ID       int       `json:"id"`
	User1    User      `json:"user1"`
	User2    User      `json:"user2"`
	Messages []Message `json:"messages"`
}

type Message struct {
	ForUser  string `json:"forUser"`
	FromUser string `json:"fromUser"`
	Body     string `json:"body"`
}

type User struct {
	ID   string `json:"id"`
	Conn *websocket.Conn
}

func (u *User) Add(r *http.Request, conn *websocket.Conn) (err error) {
	c, err := r.Cookie("user_id")
	if err != nil {
		return err
	}
	u.ID = c.Value
	// log.Println("user id ", u.ID)
	u.Conn = conn
	UserList = append(UserList, *u)
	return
}

func (m *Message) SendJSON() {
	for i, u := range UserList {
		if u.ID == m.ForUser || u.ID == m.FromUser {
			if err := u.Conn.WriteJSON(&m); err != nil {
				UserList = append(UserList[:i], UserList[i+1:]...)
			}
		}
	}
	log.Println(UserList)
}
