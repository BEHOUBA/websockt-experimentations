package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	ChatRoomList []ChatRoom
	UserList     []User
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", WShandler)

	http.ListenAndServe(":8080", nil)
}

func WShandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	var user User
	user.Add(r, conn)

	log.Println("new connection", user.ID)

	for {

		var msg Message

		// messageType, p, err := conn.ReadMessage()
		// if err != nil {
		// 	log.Println(err)
		// }
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("ERROR MESSAGE", err, &conn)
			removeUser(conn)
			break
		}
		msg.SendJSON()
	}
}

func removeUser(conn *websocket.Conn) {
	for i, u := range UserList {
		if u.Conn == conn {
			log.Println("match")
			UserList = append(UserList[:i], UserList[i+1:]...)
		}
	}
}
