package api

import (
	"fmt"
	"net/http"
)

// PingHandler checks if the API works
func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}
