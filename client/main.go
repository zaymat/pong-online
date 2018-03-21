package main

import (
	"flag"
	"log"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/veandco/go-sdl2/sdl"
)

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

func handleMessage(state *chan Map, c *websocket.Conn, player int) {
	var msg Map
	for {
		// Read message from websocket
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("error: ", err)
			return
		}

		// Check if a player wins
		if msg.Winner != 0 {
			if msg.Winner == player {
				log.Println("You win")
			} else {
				log.Println("You lose")
			}
		} else {
			// Send the message to the drawMAp routine
			*state <- msg
		}
	}
}

func drawMap(s *chan Map, window *sdl.Window) {
	// Setup drawing surface
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	// Wipe the screen
	surface.FillRect(nil, 0)

	for {
		// Read a message from the chan
		msg := <-*s

		// Wipe the screen
		surface.FillRect(nil, 0)

		// Create the 2 rackets and the ball
		player1 := sdl.Rect{0, int32(msg.Player1), 2, 28}
		player2 := sdl.Rect{509, int32(msg.Player2), 2, 28}

		var ball sdl.Rect
		if msg.Speed.Y < 0 {
			ball = sdl.Rect{int32(msg.Ball.X), int32(msg.Ball.Y) - 1, 7, 7}
		} else {
			ball = sdl.Rect{int32(msg.Ball.X), int32(msg.Ball.Y), 7, 7}
		}

		// Fill the rectangles in white
		surface.FillRect(&player1, 0xffffff00)
		surface.FillRect(&player2, 0xffffff00)
		surface.FillRect(&ball, 0xffffff00)

		// Update graphics
		window.UpdateSurface()
	}
}

func main() {

	// Parse flags
	host := flag.String("host", "localhost", "server url")
	port := flag.Int("port", 8080, "server port")

	flag.Parse()

	hostname := *host + ":" + strconv.Itoa(*port)
	// Connect to the websocket
	u := url.URL{Scheme: "ws", Host: hostname, Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Create the state chan to handle state changes from server
	state := make(chan Map)

	// Read the first message (player ID)
	_, message, err := c.ReadMessage()
	player, _ := strconv.Atoi(string(message))
	log.Println(player)

	// Init SDL for graphics
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// Create the window (size 512*256)
	window, err := sdl.CreateWindow("Pong client", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		512, 256, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// Go routine to handle state change and graphical interface
	go handleMessage(&state, c, player)
	go drawMap(&state, window)

	// Event loop
	running := true
	for running {
		var e Event
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			// Checl whether the window is closed by the user
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				switch t.Keysym.Scancode {
				// Check if the key pressed is UP
				case sdl.SCANCODE_UP:
					e.Event = "up"
					e.Player = player
					// Send up event to the server
					c.WriteJSON(&e)
				// Check if the key pressed is DOWN
				case sdl.SCANCODE_DOWN:
					e.Event = "down"
					e.Player = player
					// Send down event to the server
					c.WriteJSON(&e)
				case sdl.SCANCODE_S:
					e.Event = "start"
					e.Player = player
					// Send start event to the server
					c.WriteJSON(&e)
				case sdl.SCANCODE_R:
					e.Event = "reset"
					e.Player = player
					// Send reset event to the server
					c.WriteJSON(&e)
				case sdl.SCANCODE_P:
					e.Event = "stop"
					e.Player = player
					// Send pause event to the server
					c.WriteJSON(&e)
				}
			}
		}
	}
}
