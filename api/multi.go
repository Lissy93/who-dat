package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/lissy93/who-dat/lib"
)

// MultiHandler handles Whois requests for multiple domains
func MultiHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure it's a GET request
	if r.Method != http.MethodGet {
		http.Error(w, "Please use a GET request", http.StatusMethodNotAllowed)
		return
	}

	// Extract domains from the query parameter
	domainsQuery := r.URL.Query().Get("domains")
	if domainsQuery == "" {
		http.Error(w, "No domains specified", http.StatusBadRequest)
		return
	}
	domains := strings.Split(domainsQuery, ",")

	// Set up a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Get Whois data for all domains
	allWhois, err := lib.GetMultiWhois(ctx, domains)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert Whois data to JSON and send the response
	respondWithJSON(w, http.StatusOK, allWhois)
}

// respondWithError sends an error response in JSON format
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sends a response in JSON format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
