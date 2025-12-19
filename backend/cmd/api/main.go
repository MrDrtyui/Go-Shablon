package main

import (
	"log"

	"app/internal/app/api"
)

func main() {
	// TODO: init app
	if err := api.Run(); err != nil {
		log.Fatal("Failed to run server http", err.Error())
	}
}
