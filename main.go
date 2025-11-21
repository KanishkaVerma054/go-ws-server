package main

import (
	"KanishkaVerma054/go-ws-server/wsServer"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", wsServer.WsHandler)
	go wsServer.HandleMessage()
	fmt.Println("Websocket server started on :4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
