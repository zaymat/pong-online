package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// handleConnection : handle connections to the websocket
func (socket *WebSocket) handleConnection() func(w http.ResponseWriter, r *http.Request) {
	counter := make([]int, 0, 2)
	id := 0
	return func(w http.ResponseWriter, r *http.Request) {
		// Manage player ids
		if len(counter) == 0 {
			counter = append(counter, 1)
			id = 1
		} else if len(counter) == 1 {
			player := counter[0]
			if player == 1 {
				counter = append(counter, 2)
				id = 2
			} else {
				counter = append(counter, 1)
				id = 1
			}
		} else {
			log.Println("Too many clients")
			return
		}

		// Start websocket
		ws, err := socket.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		// Make sure we close the connection when the function returns
		defer ws.Close()

		// Send its player id to the client
		socket.Clients[ws] = id
		msg := []byte(strconv.Itoa(id))
		ws.WriteMessage(websocket.TextMessage, msg)

		for {
			var msg Event
			err := ws.ReadJSON(&msg)

			if err != nil {
				// Delete session
				id = socket.Clients[ws]
				delete(socket.Clients, ws)
				if len(counter) == 1 {
					counter = make([]int, 0, 2)
				} else {
					counter = make([]int, 0, 2)
					counter = append(counter, 3-id)
				}

				log.Println("Connection closed")
				log.Println(err)
				break
			}

			socket.Event <- msg
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
