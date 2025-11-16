package data

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// Package-level variables for shared state
var firestoreClient *firestore.Client
var encryptionKey []byte
var credMutex = &sync.Mutex{}

// InitDB loads environment configuration, decodes and validates the AES-256 encryption key,
// and initializes the package Firestore client.
//
// InitDB reads environment variables (via .env), expects ENCRYPTION_KEY as a hex-encoded 32-byte key,
// and requires PROJECT_ID and GOOGLE_APPLICATION_CREDENTIALS for Firebase initialization.
// On success it sets the package-level encryptionKey and firestoreClient; on failure it returns a descriptive error.
func InitDB(ctx context.Context) error {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Load and validate the encryption key
	keyString := os.Getenv("ENCRYPTION_KEY")
	if keyString == "" {
		return fmt.Errorf("ENCRYPTION_KEY environment variable not set")
	}
	encryptionKey, err = hex.DecodeString(keyString)
	if err != nil {
		return fmt.Errorf("failed to decode encryption key: %v", err)
	}
	if len(encryptionKey) != 32 { // AES-256 requires a 32-byte key
		return fmt.Errorf("encryption key must be 32 bytes (64 hex characters) long, but got %d bytes", len(encryptionKey))
	}

	config := &firebase.Config{
		ProjectID: os.Getenv("PROJECT_ID"),
	}
	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, config, sa)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error getting Firestore client: %v", err)
	}
	firestoreClient = client
	return nil
}

// CloseDB closes the Firestore client and releases its resources.
// If a client is initialized, it attempts to close it and logs any error encountered.
// It is safe to call CloseDB multiple times.
func CloseDB() {
	if firestoreClient != nil {
		if err := firestoreClient.Close(); err != nil {
			log.Printf("Failed to close Firestore client: %v", err)
		}
	}
}
