package api

import (
	"fmt"
	"net/http"
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
