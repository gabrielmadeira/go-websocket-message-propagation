package main

import (
	"os"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"bufio"
	"fmt"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	u := url.URL{Scheme: "ws", Host: ":3000", Path: "/"}
	conn, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
	log.Println("Connected...")
	receiveFromServer := func() {
		for {
			_, messageBytes, _ := conn.ReadMessage()
			fmt.Println("Received Message: ", string(messageBytes))
		}
		conn.Close()
	}
	go receiveFromServer()
	for {
		text, _ := reader.ReadString('\n')
		conn.WriteMessage(websocket.BinaryMessage, []byte(text))
	}
	conn.Close()
}