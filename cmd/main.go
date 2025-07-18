package main

import (
	"log"

	"github.com/somepgs/go_final_project/pkg/db"
	"github.com/somepgs/go_final_project/pkg/server"
)

// main initializes the database and starts the server.
// It checks for the existence of the database file and creates the scheduler table if it does not exist.
func main() {
	err := db.Init() // Initialize the database
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Ensure the database is closed when the application exits
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
	server.Run() // Start the server
}
