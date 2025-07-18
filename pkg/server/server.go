package server

import (
	"log"
	"net/http"
	"os"

	"github.com/somepgs/go_final_project/pkg/api"
)

const webDir = "web"

var myServer *http.Server

func Run() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init(mux)

	myServer = &http.Server{
		Addr:     ":" + port,
		Handler:  mux,
		ErrorLog: logger,
	}

	log.Printf("Server is starting on port %s", port)
	if err := myServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
