package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var (
	users      []User
	usersMutex sync.Mutex

	logs      []LogEntry
	logsMutex sync.Mutex

	startTime    = time.Now()
	startTimeStr = startTime.In(time.FixedZone("WAT", 3600)).Format("15:04:05")
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Password     string `json:"password,omitempty"`
	Message      string `json:"message,omitempty"`
	RegisteredAt string `json:"registered_at"`
}

type LogEntry struct {
	ID        int    `json:"id"`
	Action    string `json:"action"`
	Target    string `json:"target"`
	Timestamp string `json:"timestamp"`
	Unix      int64  `json:"unix"`
}

func addLog(action, target string) {
	logsMutex.Lock()
	defer logsMutex.Unlock()
	now := time.Now().In(time.FixedZone("WAT", 3600))
	logs = append(logs, LogEntry{
		ID:        len(logs) + 1,
		Action:    action,
		Target:    target,
		Timestamp: now.Format("2006-01-02 15:04:05"),
		Unix:      now.Unix(),
	})
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
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
			username := u.Username
			users = append(users[:i], users[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully!"})
			go addLog("deleted", username)
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
			if existingUser.Username == u.Username {
				usersMutex.Unlock()
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "❌ Error: User '" + u.Username + "' already exists!",
				})
				return
			}
		}
		u.ID = len(users) + 1
		u.RegisteredAt = time.Now().In(time.FixedZone("WAT", 3600)).Format("2006-01-02")
		users = append(users, u)
		usersMutex.Unlock()

		u.Message = fmt.Sprintf("✅ Success! %s added.", u.Username)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
		go addLog("added", u.Username)
		return
	}

	if r.Method == http.MethodGet {
		usersMutex.Lock()
		json.NewEncoder(w).Encode(users)
		usersMutex.Unlock()
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	uptime := time.Since(startTime)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	usersMutex.Lock()
	totalUsers := len(users)
	usersMutex.Unlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"started_at":     startTimeStr,
		"uptime":         fmt.Sprintf("%02dh %02dm %02ds", hours, minutes, seconds),
		"uptime_secs":    int(uptime.Seconds()),
		"total_users":    totalUsers,
		"go_version":     runtime.Version(),
		"os":             runtime.GOOS,
		"arch":           runtime.GOARCH,
		"num_cpu":        runtime.NumCPU(),
		"goroutines":     runtime.NumGoroutine(),
		"alloc_mb":       fmt.Sprintf("%.2f", float64(mem.Alloc)/1024/1024),
		"total_alloc_mb": fmt.Sprintf("%.2f", float64(mem.TotalAlloc)/1024/1024),
		"sys_mb":         fmt.Sprintf("%.2f", float64(mem.Sys)/1024/1024),
	})
}

func analyticsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usersMutex.Lock()
	defer usersMutex.Unlock()

	logsMutex.Lock()
	defer logsMutex.Unlock()

	today := time.Now().In(time.FixedZone("WAT", 3600)).Format("2006-01-02")

	regPerDay := map[string]int{}
	todayCount := 0
	for _, u := range users {
		regPerDay[u.RegisteredAt]++
		if u.RegisteredAt == today {
			todayCount++
		}
	}

	deletedCount := 0
	for _, l := range logs {
		if l.Action == "deleted" {
			deletedCount++
		}
	}

	history := []map[string]interface{}{}
	for i := 6; i >= 0; i-- {
		day := time.Now().AddDate(0, 0, -i).In(time.FixedZone("WAT", 3600)).Format("2006-01-02")
		history = append(history, map[string]interface{}{
			"date":  day,
			"count": regPerDay[day],
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_users":   len(users),
		"today":         todayCount,
		"total_deleted": deletedCount,
		"history":       history,
	})
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logsMutex.Lock()
	defer logsMutex.Unlock()

	reversed := make([]LogEntry, len(logs))
	for i, l := range logs {
		reversed[len(logs)-1-i] = l
	}
	json.NewEncoder(w).Encode(reversed)
}

func openBrowser(url string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	case "windows":
		exec.Command("cmd", "/c", "start", url).Start()
	}
}

func openAllPages() {
	pages := []string{
		"http://localhost:8080",
		"http://localhost:8080/dashboard",
		"http://localhost:8080/manage",
		"http://localhost:8080/analytics",
		"http://localhost:8080/logs",
		"http://localhost:8080/system-status",
	}
	time.Sleep(500 * time.Millisecond)
	for _, page := range pages {
		openBrowser(page)
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./static/landing.html")
	})
	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	mux.HandleFunc("/manage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/manage.html")
	})
	mux.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/analytics.html")
	})
	mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/logs.html")
	})
	mux.HandleFunc("/system-status", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/status.html")
	})

	mux.HandleFunc("/api/user", cors(userHandler))
	mux.HandleFunc("DELETE /api/user/{id}", cors(deleteUserHandler))
	mux.HandleFunc("/api/status", cors(statusHandler))
	mux.HandleFunc("/api/analytics", cors(analyticsHandler))
	mux.HandleFunc("/api/logs", cors(logsHandler))

	fmt.Println("🚀 Server running at http://localhost:8080")
	fmt.Printf("⏰ Started at: %s WAT\n", startTimeStr)
	go openAllPages()
	http.ListenAndServe(":8080", mux)
}
