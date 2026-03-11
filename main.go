package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// User matches the JSON structure sent from the frontend script.js
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// userHandler processes the "Trigger API Call" button click
func userHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests for this endpoint
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	// Decode the JSON sent by the browser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Logic: Let's modify the name slightly to show the server processed it
	u.Name = "Server says: Hello, " + u.Name

	// Send the JSON back to the frontend
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func main() {
	// Create the router
	mux := http.NewServeMux()

	// 1. Serve static files (CSS, JS) from the /static/ directory
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// 2. Serve the index.html file at the root URL (http://localhost:8080)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// 3. The API endpoint the button in index.html calls
	mux.HandleFunc("/api/user", userHandler)

	// 4. Configure the server with timeouts for safety
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	fmt.Println("Modern Go Server is running at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
