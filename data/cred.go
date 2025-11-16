package data

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	docRef, err := findSiteDocument(ctx, site)
	if err != nil {
		return err // Will be a "not found" error if neither version exists
	}

	credMutex.Lock()
	defer credMutex.Unlock()

	// Use the ID of the document that was actually found
	return updateLocked(ctx, docRef.ID, username, password)
}

// Show retrieves all credentials from Firestore, decrypts their passwords, and returns them as a slice.
// For each document, it attempts to decrypt the password. On decryption failure, the error is logged,
// and the password for that entry is set to "[DECRYPTION FAILED]".
// It returns an error if iterating through the collection fails.
func Show(ctx context.Context) ([]SiteCredentials, error) {
	var results []SiteCredentials
	iter := firestoreClient.Collection("credentials").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate credentials: %v", err)
		}

		var creds Credentials
		doc.DataTo(&creds)

		decryptedPassword, err := decrypt(creds.Password)
		if err != nil {
			log.Printf("Failed to decrypt password for site %s: %v. Omitting password.", doc.Ref.ID, err)
			decryptedPassword = "[DECRYPTION FAILED]"
		}
		results = append(results, SiteCredentials{
			Site:     doc.Ref.ID,
			Username: creds.Username,
			Password: decryptedPassword,
		})
	}
	return results, nil
}

// Retrieve fetches credentials for the given site from Firestore and decrypts the stored password.
// It attempts to find the site with and without a ".com" suffix.
// It returns the Credentials with the decrypted Password and true if a document is found and decryption succeeds; otherwise it returns an empty Credentials and false.
func Retrieve(ctx context.Context, site string) (Credentials, bool) {
	docRef, err := findSiteDocument(ctx, site)
	if err != nil {
		return Credentials{}, false // Document not found
	}

	doc, err := docRef.Get(ctx)
	if err != nil {
		log.Printf("Failed to get document snapshot for site %s: %v", docRef.ID, err)
		return Credentials{}, false
	}

	var creds Credentials
	if err := doc.DataTo(&creds); err != nil {
		log.Printf("Failed to parse data for site %s: %v", doc.Ref.ID, err)
		return Credentials{}, false
	}

	decryptedPassword, err := decrypt(creds.Password)
	if err != nil {
		log.Printf("Failed to decrypt password for site %s: %v", doc.Ref.ID, err)
		return Credentials{}, false
	}
	creds.Password = decryptedPassword
	return creds, true
}

// Delete removes the credentials document for the named site from Firestore.
// It acquires an internal mutex to serialize access and returns true on success, false on failure.
func Delete(ctx context.Context, site string) bool {
	docRef, err := findSiteDocument(ctx, site)
	if err != nil {
		return false // Site not found
	}

	credMutex.Lock()
	defer credMutex.Unlock()

	_, err = docRef.Delete(ctx)
	if err != nil {
		// We can log this error if needed, but for the function signature we just return false
		fmt.Printf("Failed to delete site %s: %v\n", site, err)
		return false
	}
	fmt.Println("Credentials deleted successfully!")
	return true
}
