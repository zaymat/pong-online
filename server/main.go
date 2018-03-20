package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Score : Store the score
type Score struct {
	Player1 int
	Player2 int
}

// Pos : store the position of the ball
type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func main() {
	// Routing part
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	// Create websocket
	var ws WebSocket
	ws.Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws.Clients = make(map[*websocket.Conn]int)
	ws.Broadcast = make(chan State)
	ws.Event = make(chan Event)

	// Handle connection to websocket
	conn := r.PathPrefix("/connect").Subrouter()
	handler := ws.handleConnection()
	conn.Methods("GET").HandlerFunc(handler)

	// Handle messages from Websocket
	go ws.handleMessages()

	// hangle game events
	go ws.game()

	// Serve the API
	http.Handle("/", r)
	fmt.Println("Server listening on port 8080 ...")
	http.ListenAndServe(":8080", r)
}
