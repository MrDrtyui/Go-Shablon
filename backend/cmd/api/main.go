package main

import (
	"log"

	"app/internal/app/app"
)

func main() {
	// TODO: init app
	if err := app.Run(); err != nil {
		log.Fatal("Failed to run server http", err.Error())
	}
}
