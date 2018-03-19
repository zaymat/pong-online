package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Pos : represent a position in a cathesian coordinate system
type Pos struct {
	X int
	Y int
}

// Map : represent the state of the map
type Map struct {
	Player1 int  `json:"player1"` // Player1 racket left high corner position (28*2)
	Player2 int  `json:"player2"` // Player2 racket left high corner position (28*2)
	Ball    Pos  `json:"ball"`    // Ball center position (7*7)
	Running bool `json:"running"` // Check whether the game is running
}

func handleMessage(state *chan Map, c *websocket.Conn) {
	var msg Map
	for {
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("error: ", err)
			return
		}
		*state <- msg
	}
}

func drawMap(s *chan Map) {
	for {
		msg := <-*s
		log.Println(msg)
	}
}

func main() {

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	state := make(chan Map)

	_, message, err := c.ReadMessage()
	player := string(message)
	log.Println(player)

	go handleMessage(&state, c)
	go drawMap(&state)

	for {

	}
}
