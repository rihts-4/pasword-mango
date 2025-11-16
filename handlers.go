package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rihts-4/pasword-mango/data"
)

func credentialsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the site from the URL path, e.g., "/credentials/google.com" -> "google.com"
	site := strings.TrimPrefix(r.URL.Path, "/credentials/")
	site = strings.TrimSuffix(site, "/")

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

		if err := data.Store(ctx, payload.Site, payload.Username, payload.Password); err != nil {
			http.Error(w, "Failed to store credentials", http.StatusInternalServerError)
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
		if !data.Delete(ctx, site) {
			http.Error(w, "Failed to delete credentials or not found", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Credentials deleted successfully.")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
