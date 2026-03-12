package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	users      []User
	usersMutex sync.Mutex // Prevents data corruption during simultaneous writes
	startTime  = time.Now().In(time.FixedZone("WAT", 3600)).Format("15:04:05")
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var u User
		json.NewDecoder(r.Body).Decode(&u)

		// Save user to the list
		usersMutex.Lock()
		users = append(users, u)
		usersMutex.Unlock()

		u.Message = fmt.Sprintf("%s added successfully!", u.Name)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
		return
	}

	// Allow frontend to GET the full list
	if r.Method == http.MethodGet {
		usersMutex.Lock()
		json.NewEncoder(w).Encode(users)
		usersMutex.Unlock()
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	mux.HandleFunc("/api/user", userHandler)
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"started_at": startTime})
	})

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
