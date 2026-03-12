# -web-server
🚀 Go System Dashboard
A modern, high-performance web server built with Golang and a sleek Glassmorphism frontend. This dashboard allows real-time user registration, live system monitoring, and persistent data handling—all powered by the Go standard library.


## ✨ Features

    - ⚡ High-Performance Backend: Built using net/http for maximum efficiency.
    - 🕒 Time Sync: Automatically synchronized to West Africa Time (WAT).
    - 🎨 Modern UI: Futuristic dark theme with neon glows and glass-blur effects.
    - 🔄 Real-time Interaction: Add users and see them appear instantly without page reloads.
    - 📡 Live Monitoring: Pulsing "Server Online" status and precise server start-time tracking.
    - 🧹 Data Management: Includes a "Clear All" function to wipe the in-memory user list.

## 📁 Project Structure

## ├── main.go            # Backend Go server & API logic
## ├── go.mod             # Go module definition
## └── static/            # Frontend assets
   ## ├── index.html     # Dashboard layout
    ## ├── style.css      # Modern Glassmorphism styling
    ## └── script.js      # Frontend logic & API calls

## 🛠️ Installation & Setup
***1. Prerequisites***

    Install Go 1.22 or higher.

## 2. Clone and Initialize
bash

# Initialize the Go module
go mod ini web-server

## 3. Run the Server
bash

# Start the backend
go run main.go

## 4. Access the Dashboard
Open your browser and navigate to:
👉 http://localhost:8080
🧪 API Endpoints
Method	Endpoint	Description

- GET	/	Serves the main Dashboard UI

- GET	/api/status	Returns server start time in WAT

- GET	/api/user	Fetches the list of all registered users

- POST	/api/user	Registers a new user to the system

## 👨‍💻 Developed By
Chekus Joseph Coder

Building fast, beautiful, and reliable systems with Go.
