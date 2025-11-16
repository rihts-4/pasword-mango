package data

import "errors"

// ErrNotFound is returned when a requested credential is not found in the database.
var ErrNotFound = errors.New("credentials not found")

// Credentials represents the data structure as it is stored in the Firestore database.
type Credentials struct {
	Username string `firestore:"username"`
	Password string `firestore:"password"`
}

// SiteCredentials extends Credentials to include the site identifier, used for API responses.
type SiteCredentials struct {
	Site     string `json:"site"`
	Username string `json:"username"`
	Password string `json:"password"`
}
