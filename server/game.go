package main

import (
	"log"
	"time"
)

// StraighLine : represented by ax+by+c=0
type Speed struct {
	vx int
	vy int
}

// State : Represent the map status
type State struct {
	Player1 int
	Player2 int
	Ball    Pos
	Speed   Speed
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
		log.Println("Move ball", s.Ball, s.Speed)
		x := s.Ball.x
		y := s.Ball.y
		vx := s.Speed.vx
		vy := s.Speed.vy

		s.Ball.x = x + vx
		s.Ball.y = y + vy

		if x < 0 || x > 512 {
			s.Speed.vx = -1 * s.Speed.vx
			if x < 0 {
				s.Ball.x = 0
			} else {
				s.Ball.x = 511
			}
		}

		if y < 0 || y > 256 {
			s.Speed.vy = -1 * s.Speed.vy
			if y < 0 {
				s.Ball.y = 0
			} else {
				s.Ball.y = 255
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (ws *WebSocket) game() {
	state := State{0, 0, Pos{0, 0}, Speed{1, 1}}
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
