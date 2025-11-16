package data

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"golang.org/x/net/publicsuffix"
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

	// If not found, try the alternative form using public suffix logic.
	// This handles multi-level TLDs like .co.uk correctly.
	if status.Code(err) == codes.NotFound {
		alternativeSite := getAlternativeSite(site)
		if alternativeSite != "" {
			docRef = firestoreClient.Collection("credentials").Doc(alternativeSite)
			if _, err = docRef.Get(ctx); err == nil {
				return docRef, nil // Found on second try
			}
		}
	}
	// Log the error if it's not a NotFound error for debugging purposes.
	if err != nil && status.Code(err) != codes.NotFound {
		log.Printf("Error during findSiteDocument for site '%s': %v", site, err)
	}
	return nil, ErrNotFound
}

// getAlternativeSite computes an alternative site name for credential lookup using public suffix awareness.
//
// For domains ending with ".com", it removes only the ".com" suffix (e.g., "example.com" → "example").
// For other domains, it replaces the effective TLD+1 (eTLD+1) with the base plus ".com",
// while preserving any subdomains (e.g., "www.example.co.uk" → "www.example.com").
//
// This approach correctly handles multi-level TLDs like .co.uk, preventing invalid alternatives
// such as "example.co.uk.com".
//
// Known limitation: This function may produce unexpected results for:
//   - Invalid or malformed domain names
//   - Domains with unusual structures that don't follow standard conventions
//   - Private/internal TLDs not in the public suffix list
//
// If the eTLD+1 cannot be determined, the function returns an empty string.
func getAlternativeSite(site string) string {
	// Handle the simple case: remove ".com" if present
	if strings.HasSuffix(site, ".com") {
		return strings.TrimSuffix(site, ".com")
	}

	// For non-.com domains, use public suffix logic to construct the alternative
	eTLDPlus1, err := publicsuffix.EffectiveTLDPlusOne(site)
	if err != nil {
		// If we can't determine the eTLD+1, fall back to simple append
		// This handles edge cases where the domain might not be in the public suffix list
		return site + ".com"
	}

	// Get the TLD (public suffix)
	eTLD, icann := publicsuffix.PublicSuffix(site)
	if !icann || eTLD == "" {
		// Fallback if we can't get the TLD or it's not an ICANN-managed domain
		return site + ".com"
	}

	// Get the base domain by removing the TLD from eTLD+1
	baseDomain := strings.TrimSuffix(eTLDPlus1, "."+eTLD)

	// If the site is exactly the eTLD+1 (no subdomains), return base.com
	if site == eTLDPlus1 {
		return baseDomain + ".com"
	}

	// Otherwise, we have subdomains. Replace the eTLD+1 portion with base + ".com"
	// Get the subdomain prefix (everything before the eTLD+1)
	subdomain := strings.TrimSuffix(site, "."+eTLDPlus1)

	// Construct the alternative: subdomain.base.com
	if subdomain != "" {
		return subdomain + "." + baseDomain + ".com"
	}

	return baseDomain + ".com"
}
