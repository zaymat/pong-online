package main

import (
	"log"

	"github.com/gorilla/websocket"
)

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
