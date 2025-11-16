package main

import (
	"context"
	"net/http"
)

// main configures HTTP handlers for the /credentials endpoints and starts an HTTP server listening on port 8080, managing the server lifecycle including graceful shutdown.
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