package main

import (
	"TurboGin/internal/wire"
	"log"
)

func main() {
	// Initialize the application using Wire dependency injection
	server, cleanup, err := wire.InitApp()
	if err != nil {

		// Since logger might not be initialized, we'll use zap's global logger
		log.Fatal(err)
	}
	// Ensure cleanup runs when the application exits
	defer cleanup()

	// Start the server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
