package main

import (
	"log"
	"net/http"

	"agones.dev/agones/sdks/go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	// Routing part
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	// Create Agones SDK
	sdk, err := sdk.NewSDK()
	if err != nil {
		log.Println("Error occur when initializing Agones SDK : ", err)
		return
	}

	err = sdk.Ready()
	if err != nil {
		log.Println("Cannot send ready state to Agones : ", err)
	}

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
	conn.Methods("GET").HandlerFunc(ws.handleConnection())

	// Handle messages from Websocket
	go ws.handleMessages()

	// hangle game events
	go ws.game(sdk)

	// health check for Agones
	go healthCheck(sdk)

	// Serve the API
	http.Handle("/", r)
	log.Println("Server listening on port 8081 ...")
	http.ListenAndServe(":8081", r)
}
