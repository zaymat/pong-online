package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// State : Represent the map status
type State struct {
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
	Ball    string `json:"ball"`
}

// Event : Represent the event sent by the client
type Event struct {
	Player int    // ID of the player
	Event  string // Event : down or up
}

// WebSocket : Create Websocket structure
type WebSocket struct {
	Broadcast chan State              // Broadcasting channel
	Event     chan Event              // Event channel
	Clients   map[*websocket.Conn]int // connected clients : int is 0 (if client is disconnected), 1 or 2.
	Upgrader  websocket.Upgrader
}

// handleConnection : handle connections to the websocket
func (socket *WebSocket) handleConnection() func(w http.ResponseWriter, r *http.Request) {
	counter := 0
	return func(w http.ResponseWriter, r *http.Request) {
		counter++
		if counter > 2 {
			log.Println("Too many clients")
		} else {
			ws, err := socket.Upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Fatal(err)
			}
			// Make sure we close the connection when the function returns
			defer ws.Close()

			socket.Clients[ws] = counter

			for {
				if socket.Clients[ws] == 0 {
					delete(socket.Clients, ws)
				}
			}
		}
	}
}

// handleMessages : handle messages sent to the websocket
func (socket *WebSocket) handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-socket.Broadcast

		// Send it out to every client that is currently connected
		for client := range socket.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(socket.Clients, client)
			}
		}
	}
}
