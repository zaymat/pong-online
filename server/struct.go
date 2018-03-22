// This file gather all the structures of the main package
package main

import "github.com/gorilla/websocket"

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

// Score : Store the score
type Score struct {
	Player1 int
	Player2 int
}

// Pos : store the position of the ball
type Pos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Speed : represented by ax+by+c=0
type Speed struct {
	Vx float64 `json:"vx"`
	Vy float64 `json:"vy"`
}

// State : Represent the map status
type State struct {
	Player1 int   `json:"player1"` // Player1 racket left high corner position (28*2)
	Player2 int   `json:"player2"` // Player2 racket left high corner position (28*2)
	Ball    Pos   `json:"ball"`    // Ball center position (7*7)
	Speed   Speed `json:"speed"`   // Speed vector of the ball
	Running bool  `json:"running"` // Check whether the game is running
	Winner  int   `json:"winner"`  // Is winner is different from 0, it represents the winner id
}
