package main

import (
	"flag"
	"log"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/veandco/go-sdl2/sdl"
)

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
	state := make(chan Board)

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
