package main

import (
	"net/http"
	"log"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}

	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		var result string
		for _, v := range string(message) {
			result = string(v) + result
		}

		if err := ws.WriteMessage(mt, []byte(result)); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
