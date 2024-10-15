package main

import (
	"dewu/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Error running the application: %v", err)
	}
}
