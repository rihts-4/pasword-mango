package data

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Credentials struct {
	Username string `firestore:"username"`
	Password string `firestore:"password"`
}

var firestoreClient *firestore.Client
var encryptionKey []byte
var credMutex = &sync.Mutex{}

// InitDB loads environment configuration, decodes and validates the AES-256 encryption key, and initializes the package Firestore client.
//
// InitDB reads environment variables (via .env), expects ENCRYPTION_KEY as a hex-encoded 32-byte key, and requires PROJECT_ID and GOOGLE_APPLICATION_CREDENTIALS for Firebase initialization. On success it sets the package-level encryptionKey and firestoreClient; on failure it returns a descriptive error.
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

// Store saves credentials for the given site in Firestore, encrypting the password
// before writing and updating existing entries when present.
//
// Store acquires a package-level mutex to serialize write operations. If a document
// for the site already exists, it updates that document; otherwise it creates a new one.
// It returns an error if checking existence, encrypting the password, or writing to
// Firestore fails.
func Store(ctx context.Context, site string, username string, password string) error {
	credMutex.Lock()
	defer credMutex.Unlock()

	// Check if credentials for the site already exist by attempting to get the document.
	// This single read operation replaces the separate documentExists check.
	doc, err := firestoreClient.Collection("credentials").Doc(site).Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return fmt.Errorf("failed to check for existing credentials for site %s: %v", site, err)
	}

	if doc != nil && doc.Exists() {
		fmt.Printf("Credentials for '%s' already exist. Updating credentials.\n", site)
		// Call updateLocked without acquiring the mutex since we already hold it
		return updateLocked(ctx, site, username, password)
	}

	encryptedPassword, err := encrypt(password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password for site %s: %v", site, err)
	}

	creds := Credentials{
		Username: username,
		Password: encryptedPassword,
	}
	_, err = firestoreClient.Collection("credentials").Doc(site).Set(ctx, creds)
	if err != nil {
		return fmt.Errorf("failed adding credential for site %s: %v", site, err)
	}
	fmt.Println("Credentials stored successfully!")
	return nil
}

// Update updates the stored credentials for the named site with the provided username and password.
// It acquires an internal lock and delegates to the internal update implementation.
// Returns an error if encryption or the Firestore write fails.
func Update(ctx context.Context, site string, username string, password string) error {
	credMutex.Lock()
	defer credMutex.Unlock()

	return updateLocked(ctx, site, username, password)
}

// updateLocked updates credentials for a site without acquiring the mutex.
// updateLocked encrypts the provided password and writes the credential document for the specified site to Firestore.
// The caller must hold credMutex.
// It returns an error if encryption fails or if the Firestore write operation fails.
func updateLocked(ctx context.Context, site string, username string, password string) error {
	encryptedPassword, err := encrypt(password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password for site %s: %v", site, err)
	}

	creds := Credentials{
		Username: username,
		Password: encryptedPassword,
	}
	_, err = firestoreClient.Collection("credentials").Doc(site).Set(ctx, creds)
	if err != nil {
		return fmt.Errorf("failed updating credential for site %s: %v", site, err)
	}
	fmt.Println("Credentials updated successfully!")
	return nil
}

// Show lists all credentials stored in Firestore and prints each site ID, username, and decrypted password.
// For each document in the "credentials" collection it attempts to decrypt the stored password; on decryption failure
// it logs the error and prints "[DECRYPTION FAILED]". Any iteration error other than end-of-collection is logged and
// causes the program to terminate.
func Show(ctx context.Context) {
	iter := firestoreClient.Collection("credentials").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		var creds Credentials
		doc.DataTo(&creds)

		decryptedPassword, err := decrypt(creds.Password)
		if err != nil {
			log.Printf("Failed to decrypt password for site %s: %v. Showing encrypted.", doc.Ref.ID, err)
			decryptedPassword = "[DECRYPTION FAILED]"
		}
		fmt.Printf("Site: %s, Username: %s, Password: %s\n", doc.Ref.ID, creds.Username, decryptedPassword)
	}
}

// Retrieve fetches credentials for the given site from Firestore and decrypts the stored password.
// It returns the Credentials with the decrypted Password and true if the document exists and decryption succeeds; otherwise it returns an empty Credentials and false.
func Retrieve(ctx context.Context, site string) (Credentials, bool) {
	doc, err := firestoreClient.Collection("credentials").Doc(site).Get(ctx)
	if err != nil {
		// Document does not exist
		return Credentials{}, false
	}

	var creds Credentials
	doc.DataTo(&creds)

	decryptedPassword, err := decrypt(creds.Password)
	if err != nil {
		log.Printf("Failed to decrypt password for site %s: %v", site, err)
		return Credentials{}, false // Or handle error more gracefully
	}
	creds.Password = decryptedPassword
	return creds, true
}

// Delete removes the credentials document for the named site from Firestore.
// It acquires an internal mutex to serialize access and returns true on success, false on failure.
func Delete(ctx context.Context, site string) bool {
	credMutex.Lock()
	defer credMutex.Unlock()

	_, err := firestoreClient.Collection("credentials").Doc(site).Delete(ctx)
	if err != nil {
		// We can log this error if needed, but for the function signature we just return false
		fmt.Printf("Failed to delete site %s: %v\n", site, err)
		return false
	}
	fmt.Println("Credentials deleted successfully!")
	return true
}

// encrypt encrypts data using AES-GCM.
func encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

// decrypt decrypts a hex-encoded AES-GCM ciphertext using the package's encryptionKey and returns the plaintext string.
// It returns an error if the input is not valid hex, the cipher/GCM cannot be initialized with the configured key, the ciphertext is too short to contain a nonce, or AEAD authentication fails.
func decrypt(ciphertextHex string) (string, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return string(plaintext), err
}