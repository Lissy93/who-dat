package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/lissy93/who-dat/lib"
)

// MainHandler handles Whois info for a single domain
func MainHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure it's a GET request
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Please use a GET request")
		return
	}

	// Check authentication if AUTH_KEY is set
	authKey := os.Getenv("AUTH_KEY")
	if authKey != "" {
		providedKey := r.Header.Get("Authorization")
		if providedKey == "" {
			http.Error(w, "API key required. See: github.com/lissy93/who-dat", http.StatusForbidden)
			return
		}

		// Remove "Bearer " prefix if present
		if len(providedKey) > 7 && providedKey[:7] == "Bearer " {
			providedKey = providedKey[7:]
		}

		if providedKey != authKey {
			http.Error(w, "Invalid API key. See: github.com/lissy93/who-dat", http.StatusForbidden)
			return
		}
	}

	// Extract domain from URL path
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		http.Error(w, "Domain not specified", http.StatusBadRequest)
		return
	}

	// Get Whois data
	whois, err := lib.GetWhois(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert Whois data to JSON
	jsonWhois, err := jsoniter.Marshal(whois)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	w.Write(jsonWhois)
}
