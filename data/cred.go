package data

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Credentials struct {
	Username string `firestore:"username"`
	Password string `firestore:"password"`
}

var firestoreClient *firestore.Client

func InitDB(ctx context.Context) error {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
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
	// Check if credentials for the site already exist.
	if _, exists := Retrieve(ctx, site); exists {
		fmt.Printf("Credentials for '%s' already exist. Calling Update instead.\n", site)
		return Update(ctx, site, username, password)
	}

	creds := Credentials{
		Username: username,
		Password: password,
	}
	_, err := firestoreClient.Collection("credentials").Doc(site).Set(ctx, creds)
	if err != nil {
		return fmt.Errorf("failed adding credential for site %s: %v", site, err)
	}
	fmt.Println("Credentials stored successfully!")
	return nil
}

func Update(ctx context.Context, site string, username string, password string) error {
	creds := Credentials{
		Username: username,
		Password: password,
	}
	_, err := firestoreClient.Collection("credentials").Doc(site).Set(ctx, creds)
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
		fmt.Printf("Site: %s, Username: %s, Password: %s\n", doc.Ref.ID, creds.Username, creds.Password)
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
	return creds, true
}

func Delete(ctx context.Context, site string) bool {
	_, err := firestoreClient.Collection("credentials").Doc(site).Delete(ctx)
	if err != nil {
		// We can log this error if needed, but for the function signature we just return false
		fmt.Printf("Failed to delete site %s: %v\n", site, err)
		return false
	}
	fmt.Println("Credentials deleted successfully!")
	return true
}
