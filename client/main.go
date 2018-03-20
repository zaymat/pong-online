package main

import (
	"log"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/veandco/go-sdl2/sdl"
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

// Event : Represent the event sent by the client
type Event struct {
	Player int    // ID of the player
	Event  string // Event : down or up
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

func drawMap(s *chan Map, window *sdl.Window) {
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	for {
		msg := <-*s

		surface.FillRect(nil, 0)
		player1 := sdl.Rect{0, int32(msg.Player1), 2, 28}
		player2 := sdl.Rect{509, int32(msg.Player2), 2, 28}
		ball := sdl.Rect{int32(msg.Ball.X), int32(msg.Ball.Y), 7, 7}
		surface.FillRect(&player1, 0xffffff00)
		surface.FillRect(&player2, 0xffffff00)
		surface.FillRect(&ball, 0xffffff00)
		window.UpdateSurface()
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
	player, _ := strconv.Atoi(string(message))
	log.Println(player)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		512, 256, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	go handleMessage(&state, c)
	go drawMap(&state, window)

	running := true
	for running {
		var e Event
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				switch t.Keysym.Scancode {
				case sdl.SCANCODE_A:
					e.Event = "up"
					e.Player = player
					c.WriteJSON(&e)
				case sdl.SCANCODE_Q:
					e.Event = "down"
					e.Player = player
					c.WriteJSON(&e)
				}

			}

		}
	}
}
