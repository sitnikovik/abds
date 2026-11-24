package main

import (
	"log"

	"energosbyt/internal/cmd/streamer"
)

func main() {
	err := streamer.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
