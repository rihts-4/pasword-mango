package main

import (
	"context"
	"net/http"
)

func main() {
	// Set up HTTP handlers
	http.HandleFunc("/credentials", credentialsHandler)
	http.HandleFunc("/credentials/", credentialsHandler)

	// Create and configure the HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil, // Use DefaultServeMux
	}

	// Start the server and handle graceful shutdown
	run(context.Background(), srv)
}
