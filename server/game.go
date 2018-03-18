package main

import (
	"log"
	"time"
)

// StraighLine : represented by ax+by+c=0
type StraighLine struct {
	a int
	b int
	c int
}

// State : Represent the map status
type State struct {
	Player1   int
	Player2   int
	Ball      Pos
	Direction StraighLine
}

// up : Update the state in case of a up command
func (s *State) up(e Event) {
	log.Println("up, Player : ", e.Player)
}

// down : Update the state in case of a down command
func (s *State) down(e Event) {
	log.Println("down, Player : ", e.Player)
}

// moveBall : move the ball on the map
func (s *State) moveBall() {
	for {
		log.Println("Move ball")

		time.Sleep(1 * time.Second)
	}
}

func (ws *WebSocket) game() {
	state := State{256, 256, Pos{0, 0}, StraighLine{1, 1, 0}}
	go state.moveBall()
	for {
		event := <-ws.Event
		switch event.Event {
		case "up":
			state.up(event)
		case "down":
			state.down(event)
		default:
			log.Println("Unknown command")
		}
	}

}
