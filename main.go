package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []byte)	// Broadcast channel
var mutex = &sync.Mutex{}	// Protect clients map

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error Upgrading:", err)
		return
	}
	defer conn.Close()
	
	// go handleConnection(conn)

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		broadcast <- message
	}
}

// func handleConnection(conn *websocket.Conn) {
// 	// Listen for incoming messages
// 	for {
// 		// Read message from client
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("Error reading message:", err)
// 			break
// 		}
// 		fmt.Printf("Received: %s\n", message)
// 		// Echo the message back to the client
// 		if err := conn.WriteMessage(websocket.TextMessage, message);
// 		err != nil {
// 			fmt.Println("Error writing message:", err)
// 			break
// 		}
// 	}
// }

func handleMessage() {
	for {
		// Grab the next message from the broadcast channel
		message := <-broadcast

		// Send the message to all connected clients
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client .Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	go handleMessage()
	fmt.Println("Websocket server started on :4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}