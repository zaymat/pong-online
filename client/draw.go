package main

import "github.com/veandco/go-sdl2/sdl"

func drawMap(s *chan Board, window *sdl.Window) {
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
		player1 := sdl.Rect{X: 0, Y: int32(msg.Player1), W: 2, H: 28}
		player2 := sdl.Rect{X: 509, Y: int32(msg.Player2), W: 2, H: 28}

		var ball sdl.Rect
		if msg.Speed.Y < 0 {
			ball = sdl.Rect{X: int32(msg.Ball.X), Y: int32(msg.Ball.Y) - 1, W: 7, H: 7}
		} else {
			ball = sdl.Rect{X: int32(msg.Ball.X), Y: int32(msg.Ball.Y), W: 7, H: 7}
		}

		// Fill the rectangles in white
		surface.FillRect(&player1, 0xffffff00)
		surface.FillRect(&player2, 0xffffff00)
		surface.FillRect(&ball, 0xffffff00)

		// Update graphics
		window.UpdateSurface()
	}
}
