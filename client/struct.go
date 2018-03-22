package main

// Pos : represent a position in a cathesian coordinate system
type Pos struct {
	X float64
	Y float64
}

// Map : represent the state of the map
type Map struct {
	Player1 int  `json:"player1"` // Player1 racket left high corner position (28*2)
	Player2 int  `json:"player2"` // Player2 racket left high corner position (28*2)
	Ball    Pos  `json:"ball"`    // Ball center position (7*7)
	Running bool `json:"running"` // Check whether the game is running
	Speed   Pos  `json:"speed"`   // Ball speed
	Winner  int  `json:"winner"`  // Winner id
}

// Event : Represent the event sent by the client
type Event struct {
	Player int    // ID of the player
	Event  string // Event : down or up
}
