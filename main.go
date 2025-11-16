package main

import (
	"context"
	"net/http"
)

// main configures HTTP handlers for the /credentials endpoints and starts an HTTP server listening on port 8080, managing the server lifecycle including graceful shutdown.
func main() {
	mux := http.NewServeMux()

	// Set up HTTP handlers
	// Wrap the handler with the logging middleware
	loggedCredentialsHandler := loggingMiddleware(http.HandlerFunc(credentialsHandler))
	mux.Handle("/credentials", loggedCredentialsHandler)
	mux.Handle("/credentials/", loggedCredentialsHandler)

	// Create and configure the HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server and handle graceful shutdown
	run(context.Background(), srv)
}
