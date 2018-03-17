package main

import (
	"fmt"
	"log"
	"net/http"

	"./handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Message : Create message structure. Represent the map status
type Message struct {
	// Player1 string `json:"player1"`
	// Player2 string `json:"player2"`
	// Ball    string `json:"ball"`
	Message string
}

// WebSocket : Create Websocket structure
type WebSocket struct {
	Broadcast chan Message             // Broadcasting channel
	Clients   map[*websocket.Conn]bool // connected clients
	Upgrader  websocket.Upgrader
}

// handleConnection : handle connections to the websocket
func (socket *WebSocket) handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := socket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	socket.Clients[ws] = true

	for {
		if socket.Clients[ws] == false {
			delete(socket.Clients, ws)
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

func main() {
	// Routing part
	r := mux.NewRouter()
	r.HandleFunc("/", home.HomeHandler)

	// Create websocket
	var ws WebSocket
	ws.Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws.Clients = make(map[*websocket.Conn]bool)
	ws.Broadcast = make(chan Message)

	// Handle connection to websocket
	conn := r.PathPrefix("/conn").Subrouter()
	conn.Path("/connect").HandlerFunc(ws.handleConnection)

	// Handle messages from Websocket
	go ws.handleMessages()

	http.Handle("/", r)
	fmt.Println("Server listening on port 8080 ...")
	http.ListenAndServe(":8080", r)
}
