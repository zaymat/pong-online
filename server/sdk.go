package main

import (
	"log"
	"time"

	"agones.dev/agones/sdks/go"
)

func healthCheck(sdk *sdk.SDK) {
	for {
		err := sdk.Health()
		if err != nil {
			log.Println("Connection problem")
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}
