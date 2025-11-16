package data

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// updateLocked updates credentials for a site without acquiring the mutex.
// updateLocked encrypts the provided password and writes the credential document for the specified site to Firestore.
// The caller must hold credMutex.
// updateLocked updates credentials for a site; the caller must hold credMutex.
// It encrypts the provided password and stores the username and encrypted password
// in the Firestore "credentials" collection using the site as the document ID.
// On success it prints a confirmation message. It returns an error if encryption
// fails or if writing the credentials to Firestore fails.
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

// findSiteDocument attempts to find a document for a given site, trying both with and without a ".com" suffix.
// findSiteDocument locates the Firestore document reference for a site's credentials.
// It first attempts the provided site string; if that document is not found it will
// try the alternative form obtained by adding or removing a trailing ".com". If a
// document is found it returns its reference. Non-NotFound errors encountered while
// querying are logged for debugging. If no document is found the function returns an error
// indicating credentials for the site were not found.
func findSiteDocument(ctx context.Context, site string) (*firestore.DocumentRef, error) {
	// First attempt: try the site name as provided.
	docRef := firestoreClient.Collection("credentials").Doc(site)
	_, err := docRef.Get(ctx)

	if err == nil {
		return docRef, nil // Found on first try
	}

	// If not found, try the alternative form.
	if status.Code(err) == codes.NotFound {
		var alternativeSite string
		if strings.HasSuffix(site, ".com") {
			alternativeSite = strings.TrimSuffix(site, ".com")
		} else {
			alternativeSite = site + ".com"
		}
		docRef = firestoreClient.Collection("credentials").Doc(alternativeSite)
		if _, err = docRef.Get(ctx); err == nil {
			return docRef, nil // Found on second try
		}
	}
	// Log the error if it's not a NotFound error for debugging purposes.
	if err != nil && status.Code(err) != codes.NotFound {
		log.Printf("Error during findSiteDocument for site '%s': %v", site, err)
	}
	return nil, fmt.Errorf("credentials for site '%s' not found", site)
}