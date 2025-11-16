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

func CloseDB() {
	if firestoreClient != nil {
		if err := firestoreClient.Close(); err != nil {
			log.Printf("Failed to close Firestore client: %v", err)
		}
	}
}

func Store(ctx context.Context, site string, username string, password string) error {
	credMutex.Lock()
	defer credMutex.Unlock()

	// Check if credentials for the site already exist.
	if exists, err := documentExists(ctx, site); err == nil && exists {
		fmt.Printf("Credentials for '%s' already exist. Updating credentials.\n", site)
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

func Update(ctx context.Context, site string, username string, password string) error {
	credMutex.Lock()
	defer credMutex.Unlock()

	return updateLocked(ctx, site, username, password)
}

// updateLocked updates credentials for a site without acquiring the mutex.
// The caller must already hold credMutex.
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

// decrypt decrypts data using AES-GCM.
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

// documentExists checks if a document for a given site exists without reading its data.
func documentExists(ctx context.Context, site string) (bool, error) {
	doc, err := firestoreClient.Collection("credentials").Doc(site).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to check for document %s: %v", site, err)
	}
	return doc.Exists(), nil
}
