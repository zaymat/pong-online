package main

import (
	"log"
)

// handleEvents : handle event from Event chan
func (ws *WebSocket) handleEvents() {
	for {
		event := <-ws.Event
		switch event.Event {
		case "up":
			Up(event)
		case "down":
			Down(event)
		default:
			log.Println("Unknown event")
		}
	}
}
