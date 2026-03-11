package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// User defines the data structure for our JSON API.
// The `json:"..."` tags tell Go how to map struct fields to JSON keys.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// userHandler processes incoming JSON data sent via POST requests.
func userHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Security Check: Only allow POST methods for this specific endpoint.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decoding: Read the JSON from the request body and save it into a User struct.
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 3. Response: Set the header to JSON, send a 201 Created status,
	// and send the decoded data back to the client to confirm success.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func main() {
	// Initialize a NewServeMux (Router) to handle our URL patterns.
	mux := http.NewServeMux()

	// 4. Static Asset Handling:
	// This tells Go to look in the "./static" folder for CSS, JS, and images.
	// StripPrefix ensures the server doesn't look for a literal "/static/" folder inside "./static".
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// 5. Root Route:
	// When a user visits "http://localhost:8080/", we manually serve the index.html file.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// 6. API Route: Map our userHandler logic to the /api/user URL.
	mux.HandleFunc("/api/user", userHandler)

	// 7. Server Configuration:
	// We define custom timeouts to prevent "hanging" connections from wasting server resources.
	server := &http.Server{
		Addr:         ":8080",          // The port the server listens on
		Handler:      mux,              // The router we defined above
		ReadTimeout:  5 * time.Second,  // Max time to read the request from the client
		WriteTimeout: 10 * time.Second, // Max time to write the response back
	}

	fmt.Println("Server starting at http://localhost:8080")

	// 8. Start the Server:
	// log.Fatal will catch and print any errors if the server fails to start.
	log.Fatal(server.ListenAndServe())
}
