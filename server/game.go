package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Speed : represented by ax+by+c=0
type Speed struct {
	Vx int `json:"vx"`
	Vy int `json:"vy"`
}

// State : Represent the map status
type State struct {
	Player1 int   `json:"player1"` // Player1 racket left high corner position (28*2)
	Player2 int   `json:"player2"` // Player2 racket left high corner position (28*2)
	Ball    Pos   `json:"ball"`    // Ball center position (7*7)
	Speed   Speed `json:"speed"`   // Speed vector of the ball
	Running bool  `json:"running"` // Check whether the game is running
}

// up : Update the state in case of a up command
func (s *State) up(e Event) {
	if e.Player == 1 {
		if s.Player1 < 227 {
			s.Player1 += 6
		}
	} else {
		if s.Player2 < 227 {
			s.Player2 += 6
		}
	}
	log.Println("up, Player : ", e.Player)
}

// down : Update the state in case of a down command
func (s *State) down(e Event) {
	if e.Player == 1 {
		if s.Player1 > 0 {
			s.Player1 -= 6
		}
	} else {
		if s.Player2 > 0 {
			s.Player2 -= 6
		}
	}
	log.Println("down, Player : ", e.Player)
}

// moveBall : move the ball on the map. Return the ID of the winner
func (s *State) moveBall(ws *WebSocket) int {
	for {
		// Check if the game is still running
		if s.Running == false {
			return 0
		}
		x := s.Ball.X
		y := s.Ball.Y
		vx := s.Speed.Vx
		vy := s.Speed.Vy

		s.Ball.X = x + vx
		s.Ball.Y = y + vy

		x = s.Ball.X
		y = s.Ball.Y

		if x < 5 || x > 506 {
			if x < 5 {
				if y >= s.Player1 && y <= s.Player1+28 {
					s.Ball.X = 5
					s.Speed.Vx = -1 * s.Speed.Vx
				} else {
					s.Ball.X = 0
					return 2
				}
			} else {
				if y >= s.Player2 && y <= s.Player2+28 {
					s.Ball.X = 506
					s.Speed.Vx = -1 * s.Speed.Vx
				} else {
					s.Ball.X = 511
					return 1
				}
			}
		}

		if y < 3 || y > 253 {
			s.Speed.Vy = -1 * s.Speed.Vy
			if y < 3 {
				s.Ball.Y = 3
			} else {
				s.Ball.Y = 252
			}
		}
		ws.Broadcast <- *s
		time.Sleep(20 * time.Millisecond)
	}
}

func (ws *WebSocket) startHandler(w http.ResponseWriter, r *http.Request) {
	var e Event
	e.Player = 0
	e.Event = "start"
}

func (ws *WebSocket) stopHandler(w http.ResponseWriter, r *http.Request) {
	var e Event
	e.Player = 0
	e.Event = "stop"
}

func (ws *WebSocket) resetHandler(w http.ResponseWriter, r *http.Request) {
	var e Event
	e.Player = 0
	e.Event = "reset"
}

func (ws *WebSocket) game() {
	r := rand.New(rand.NewSource((time.Now()).Unix()))
	x := r.Intn(505) + 3
	y := r.Intn(255)
	state := State{0, 0, Pos{x, y}, Speed{1, 1}, false}

	for {
		switch event := <-ws.Event; event.Event {
		case "up":
			state.up(event)
		case "down":
			state.down(event)
		case "start":
			state.Running = true
			go state.moveBall(ws)
		case "stop":
			state.Running = false
		case "reset":
			x := r.Intn(505) + 3
			y := r.Intn(255)
			state = State{0, 0, Pos{x, y}, Speed{1, 1}, false}
		default:
			log.Println("Unknown command")
		}
	}

}
