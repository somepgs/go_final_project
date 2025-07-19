package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/somepgs/go_final_project/pkg/api"
)

const webDir = "web"

var srv *http.Server

func Run(port int, password string) {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init(mux, password)

	adr := fmt.Sprintf(":%d", port)
	srv = &http.Server{
		Addr:     adr,
		Handler:  mux,
		ErrorLog: logger,
	}

	log.Printf("Server is starting on port %d", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
