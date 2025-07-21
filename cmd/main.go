package main

import (
	"log"
	"os"
	"strconv"

	"github.com/somepgs/go_final_project/pkg/db"
	"github.com/somepgs/go_final_project/pkg/server"
)

// config holds the configuration for the application.
type config struct {
	Port     int
	DBFile   string
	Password string
}

// envOr retrieves the value of the environment variable named by key.
func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

// loadConfig loads the configuration from environment variables.
func loadConfig() config {
	port, err := strconv.Atoi(envOr("TODO_PORT", "7540")) // Default port is 7540
	if err != nil || port <= 0 || port > 65535 {
		log.Fatalf("invalid TODO_PORT: %v", err)
	}
	return config{
		Port:     port,
		DBFile:   envOr("TODO_DBFILE", "scheduler.db"), // Default database file is scheduler.db
		Password: envOr("TODO_PASSWORD", "12345"),      // Default password is 12345
	}
}

// main initializes the database and starts the server.
// It checks for the existence of the database file and creates the scheduler table if it does not exist.
func main() {
	cfg := loadConfig()

	err := db.Init(cfg.DBFile) // Initialize the database
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Ensure the database is closed when the application exits
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
	server.Run(cfg.Port, cfg.Password) // Start the server
}
