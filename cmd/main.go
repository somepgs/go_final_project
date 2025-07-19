package main

import (
	"log"
	"os"
	"strconv"

	"github.com/somepgs/go_final_project/pkg/db"
	"github.com/somepgs/go_final_project/pkg/server"
)

type config struct {
	Port     int
	DBFile   string
	Password string
}

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func mustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}

func loadConfig() config {
	port, err := strconv.Atoi(envOr("TODO_PORT", "7540"))
	if err != nil || port <= 0 || port > 65535 {
		log.Fatalf("invalid TODO_PORT: %v", err)
	}
	return config{
		Port:     port,
		DBFile:   envOr("TODO_DBFILE", "../scheduler.db"),
		Password: mustEnv("TODO_PASSWORD"),
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
