package main

import (
	"log"
)

// Up : handle up command
func Up(e Event) {
	log.Println("up")
}

// Down : handle down command
func Down(e Event) {
	log.Println("down")
}

func (ws *WebSocket) game() {

}
