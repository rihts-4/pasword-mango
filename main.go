package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rihts-4/pasword-mango/data"
)

func main() {
	ctx := context.Background()

	// Initialize the database connection from the data package
	err := data.InitDB(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Defer closing the connection until the main function exits
	defer data.CloseDB()

	fmt.Println("--- Storing new credentials for 'google.com' ---")
	err = data.Store(ctx, "google.com", "testuser", "password123")
	if err != nil {
		log.Printf("Could not store credentials: %v\n", err)
	}
	fmt.Println()

	fmt.Println("--- Showing all stored credentials ---")
	data.Show(ctx)
	fmt.Println()

	fmt.Println("--- Retrieving credentials for 'google.com' ---")
	if creds, found := data.Retrieve(ctx, "google.com"); found {
		fmt.Printf("Found credentials for google.com: Username=%s, Password=%s\n", creds.Username, creds.Password)
	} else {
		fmt.Println("Could not find credentials for google.com.")
	}
	fmt.Println()

	// fmt.Println("--- Deleting credentials for 'google.com' ---")
	// data.Delete(ctx, "google.com")
	fmt.Println()
}
