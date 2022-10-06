package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var messagesChannel = make(chan string)
var connections = make(map[*websocket.Conn]bool)

func connectionLoop(conn *websocket.Conn) {
	log.Println("New Connection")
	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			connections[conn] = false
			log.Println("Connection Closed.")
			break
		}
		messagesChannel <- string(messageBytes)
	}
	conn.Close()
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	connections[conn] = true
	go connectionLoop(conn)
}

func broadcast() {
	for {
		message := <-messagesChannel
		for conn := range connections {
			conn.WriteMessage(websocket.BinaryMessage, []byte(message))
		}
		fmt.Println(message)
	}
}

func main() {
	go broadcast()
	http.HandleFunc("/", handler)
	log.Println("Server Listening on :3000")
	http.ListenAndServe(":3000", nil)
}