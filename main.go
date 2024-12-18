package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lissy93/who-dat/api"
	"github.com/lissy93/who-dat/lib"
)

//go:embed dist/*
var staticAssets embed.FS

func main() {
	// Create a sub-directory filesystem from the embedded files
	subFS, err := fs.Sub(staticAssets, "dist")
	if err != nil {
		log.Fatal(err)
	}

	// Create a file server for the sub-directory filesystem
	embeddedServer := http.FileServer(http.FS(subFS))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")

		if path == "docs" {
			http.ServeFile(w, r, "dist/docs.html")
			return
		}

		// Serve embedded static files for the root
		if path == "" {
			embeddedServer.ServeHTTP(w, r)
			return
		}

		// Serve embedded static files if path starts with "assets"
		if strings.HasPrefix(path, "assets") {
			embeddedServer.ServeHTTP(w, r)
			return
		}

		// Wrap API handlers with auth middleware
		if path == "multi" {
			lib.AuthMiddleware(api.MultiHandler).ServeHTTP(w, r)
		} else {
			lib.AuthMiddleware(api.MainHandler).ServeHTTP(w, r)
		}
	})

	// Choose the port to start server on
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddress := fmt.Sprintf(":%s", port)

	asciiArt := `
__          ___             _____        _  ___  
\ \        / / |           |  __ \      | ||__ \ 
 \ \  /\  / /| |__   ___   | |  | | __ _| |_  ) |
  \ \/  \/ / | '_ \ / _ \  | |  | |/ _` + "`" + ` | __|/ / 
   \  /\  /  | | | | (_) | | |__| | (_| | |_|_|  
    \/  \/   |_| |_|\___/  |_____/ \__,_|\__(_)																							
`
	log.Println(asciiArt)
	log.Printf("\nWelcome to Who-Dat - WHOIS Lookup Service.\nApp up and running at %s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}
