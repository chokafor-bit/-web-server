package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Global variable to capture when the server starts
var serverStartTime = time.Now().Format("15:04:05")

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"` // Used for the success response
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// MODIFIED: This sets the message that your JavaScript will display
	u.Message = fmt.Sprintf("✅ Success! %s has been added to the system.", u.Name)

	// Keep your original logic as well
	u.Name = "Server says: Hello, " + u.Name

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func main() {
	mux := http.NewServeMux()

	// 1. Serve static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// 2. Serve the index.html file
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// 3. API endpoint for user registration
	mux.HandleFunc("/api/user", userHandler)

	// 4. NEW: API endpoint to provide the server start time
	mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"started_at": serverStartTime,
		})
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	fmt.Printf("🚀 Modern Go Server is running at http://localhost:8080\n")
	fmt.Printf("⏰ Server started at: %s\n", serverStartTime)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
