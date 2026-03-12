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
	usersMutex sync.Mutex
	startTime  = time.Now().In(time.FixedZone("WAT", 3600)).Format("15:04:05")
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	usersMutex.Lock()
	defer usersMutex.Unlock()

	for i, u := range users {
		if fmt.Sprintf("%d", u.ID) == idStr {
			users = append(users[:i], users[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully!"})
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var u User
		json.NewDecoder(r.Body).Decode(&u)

		usersMutex.Lock()
		for _, existingUser := range users {
			if existingUser.Name == u.Name {
				usersMutex.Unlock()
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "❌ Error: User '" + u.Name + "' already exists!",
				})
				return
			}
		}

		users = append(users, u)
		usersMutex.Unlock()

		u.Message = fmt.Sprintf("✅ Success! %s added.", u.Name)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
		return
	}

	if r.Method == http.MethodGet {
		usersMutex.Lock()
		json.NewEncoder(w).Encode(users)
		usersMutex.Unlock()
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// --- UPDATED PAGE ROUTES ---

	// 1. Landing Page (The first thing users see)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Ensure only exactly "/" matches this, otherwise return 404
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./static/landing.html")
	})

	// 2. Dashboard (Moved from / to /dashboard)
	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// 3. Management Page
	mux.HandleFunc("/manage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/manage.html")
	})

	// --- API ROUTES ---
	mux.HandleFunc("/api/user", userHandler)
	mux.HandleFunc("DELETE /api/user/{id}", deleteUserHandler)
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"started_at": startTime})
	})

	fmt.Println("🚀 Server running at http://localhost:8080")
	fmt.Printf("⏰ Started at: %s WAT\n", startTime)
	http.ListenAndServe(":8080", mux)
}
