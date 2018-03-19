package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

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

			// Send its player id to the client
			socket.Clients[ws] = counter
			msg := []byte(strconv.Itoa(counter))
			ws.WriteMessage(websocket.TextMessage, msg)

			for {
				var msg Event
				err := ws.ReadJSON(&msg)

				if err != nil {
					delete(socket.Clients, ws)
					counter--
					log.Println("Connection closed")
					log.Println(err)
					break
				}

				socket.Event <- msg
			}
		}
	}
}

// handleMessages : handle messages sent to the websocket
func (socket *WebSocket) handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-socket.Broadcast
		log.Println(msg)
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
