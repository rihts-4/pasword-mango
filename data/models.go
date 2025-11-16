package data

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
