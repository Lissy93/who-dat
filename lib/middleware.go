package lib

import (
	"net/http"
	"os"
)

// AuthMiddleware checks for valid API key if AUTH_KEY env var is set
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authKey := os.Getenv("AUTH_KEY")
		
		// If no auth key is set, skip authentication
		if authKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Check the API key from the Authorization header
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

		next.ServeHTTP(w, r)
	}
} 
