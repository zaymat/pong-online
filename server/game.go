package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"agones.dev/agones/sdks/go"
)

// down : Update the state in case of a down command
func (s *State) down(e Event) {
	if e.Player == 1 {
		if s.Player1 < 227 {
			s.Player1 += 6
		}
	} else {
		if s.Player2 < 227 {
			s.Player2 += 6
		}
	}
}

// up : Update the state in case of a up command
func (s *State) up(e Event) {
	if e.Player == 1 {
		if s.Player1 > 0 {
			s.Player1 -= 6
		}
	} else {
		if s.Player2 > 0 {
			s.Player2 -= 6
		}
	}
}

// moveBall : move the ball on the map. Return the ID of the winner
func (s *State) moveBall(ws *WebSocket) {
	for {
		// Check if the game is still running
		if s.Running == false {
			return
		}
		x := s.Ball.X
		y := s.Ball.Y
		vx := s.Speed.Vx
		vy := s.Speed.Vy

		// Change ball position according to speed vector
		s.Ball.X = x + vx
		s.Ball.Y = y + vy

		x = s.Ball.X
		y = s.Ball.Y

		// Check collisions
		if x < 2 || x > 502 {
			if x < 2 {
				if int(y+7) >= s.Player1 && int(y) <= s.Player1+28 {
					s.Ball.X = 2
					s.Speed.Vx = -1 * s.Speed.Vx
					s.Speed.Vy = (0.7 + 0.05*math.Abs(float64(s.Player1)+10.0-y)) * s.Speed.Vy / math.Abs(s.Speed.Vy)
				} else {
					s.Ball.X = 0
					s.Winner = 2
					return
				}
			} else {
				if int(y+7) >= s.Player2 && int(y) <= s.Player2+28 {
					s.Ball.X = 502
					s.Speed.Vx = -1 * s.Speed.Vx
					s.Speed.Vy = (0.7 + 0.05*math.Abs(float64(s.Player2)+10.0-y)) * s.Speed.Vy / math.Abs(s.Speed.Vy)
				} else {
					s.Ball.X = 511
					s.Winner = 1
					return
				}
			}
		}

		if y < 0 || y > 248 {
			s.Speed.Vy = -1 * s.Speed.Vy
			if y < 0 {
				s.Ball.Y = 0
			} else {
				s.Ball.Y = 248
			}
		}
		time.Sleep(20 * time.Millisecond)
		ws.Broadcast <- *s
	}
}

// Handle events from client
func (ws *WebSocket) game(sdk *sdk.SDK) {
	r := rand.New(rand.NewSource((time.Now()).Unix()))
	x := float64(r.Intn(505) + 3)
	y := float64(r.Intn(248))
	state := State{0, 0, Pos{x, y}, Speed{1, 1}, false, 0}
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
			err := sdk.Shutdown()
			if err != nil {
				log.Println("Unable to shutdown")
			}
		case "reset":
			x := float64(r.Intn(505) + 3)
			y := float64(r.Intn(248))
			state = State{0, 0, Pos{x, y}, Speed{1, 1}, false, 0}
		default:
			log.Println("Unknown command")
		}
		ws.Broadcast <- state
	}

}
