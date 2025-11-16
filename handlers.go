package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rihts-4/pasword-mango/data"
)

// credentialsHandler handles HTTP CRUD operations for credentials under the /credentials/ path.
//
// It supports the following methods:
// - POST: creates credentials from a JSON body containing `site`, `username`, and `password` (returns 201 on success).
// - GET: without a site lists all credentials as JSON; with a site returns that site's credentials as JSON (returns 404 if not found).
// - PUT: updates credentials for the site in the URL using a JSON body with `username` and `password` (site must be present in the path).
// - DELETE: deletes credentials for the site in the URL (site must be present in the path).
//
// The handler returns 400 for malformed requests, 500 for internal/data errors, and 405 for unsupported methods.
func credentialsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the site from the URL path, e.g., "/credentials/google.com" -> "google.com"
	path := strings.TrimPrefix(r.URL.Path, "/credentials")
	path = strings.Trim(path, "/")
	site := path // empty => list all; non-empty => specific site
	ctx := r.Context()

	switch r.Method {
	case http.MethodPost: // Create new credentials
		var payload struct {
			Site     string `json:"site"`
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Trim whitespace from all fields
		payload.Site = strings.TrimSpace(payload.Site)
		payload.Username = strings.TrimSpace(payload.Username)
		payload.Password = strings.TrimSpace(payload.Password)

		// Validate that all fields are non-empty
		if payload.Site == "" {
			http.Error(w, "Site is required and cannot be empty", http.StatusBadRequest)
			return
		}
		if payload.Username == "" {
			http.Error(w, "Username is required and cannot be empty", http.StatusBadRequest)
			return
		}
		if payload.Password == "" {
			http.Error(w, "Password is required and cannot be empty", http.StatusBadRequest)
			return
		}

		// Enforce reasonable length limits
		const maxSiteLength = 255
		const maxUsernameLength = 255
		const maxPasswordLength = 1000

		if len(payload.Site) > maxSiteLength {
			http.Error(w, fmt.Sprintf("Site must not exceed %d characters", maxSiteLength), http.StatusBadRequest)
			return
		}
		if len(payload.Username) > maxUsernameLength {
			http.Error(w, fmt.Sprintf("Username must not exceed %d characters", maxUsernameLength), http.StatusBadRequest)
			return
		}
		if len(payload.Password) > maxPasswordLength {
			http.Error(w, fmt.Sprintf("Password must not exceed %d characters", maxPasswordLength), http.StatusBadRequest)
			return
		}

		if err := data.Store(ctx, payload.Site, payload.Username, payload.Password); err != nil {
			if err == data.ErrAlreadyExists {
				http.Error(w, "Site already exists. Use PUT to update.", http.StatusConflict)
				return
			}
			http.Error(w, fmt.Sprintf("Failed to store credentials: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Credentials stored successfully.")

	case http.MethodGet:
		if site == "" { // Show all credentials
			creds, err := data.Show(ctx)
			if err != nil {
				http.Error(w, "Failed to retrieve credentials", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(creds)
		} else { // Retrieve specific credentials
			creds, found := data.Retrieve(ctx, site)
			if !found {
				http.Error(w, "Credentials not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(creds)
		}

	case http.MethodPut: // Update credentials
		if site == "" {
			http.Error(w, "Missing site in URL path for update", http.StatusBadRequest)
			return
		}
		var creds data.Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// Trim and validate username and password
		creds.Username = strings.TrimSpace(creds.Username)
		creds.Password = strings.TrimSpace(creds.Password)
		if creds.Username == "" || creds.Password == "" {
			http.Error(w, "username and password required", http.StatusBadRequest)
			return
		}
		if err := data.Update(ctx, site, creds.Username, creds.Password); err != nil {
			http.Error(w, "Failed to update credentials", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Credentials updated successfully.")

	case http.MethodDelete: // Delete credentials
		if site == "" {
			http.Error(w, "Missing site in URL path for delete", http.StatusBadRequest)
			return
		}
		if err := data.Delete(ctx, site); err != nil {
			if err == data.ErrNotFound {
				http.Error(w, "Credentials not found", http.StatusNotFound)
				return
			}
			http.Error(w, fmt.Sprintf("Failed to delete credentials: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Credentials deleted successfully.")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
