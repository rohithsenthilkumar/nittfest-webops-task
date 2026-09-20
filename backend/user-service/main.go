package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// User represents a user entity
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Hardcoded users
var users = []User{
	{ID: 1, Name: "Alice Johnson", Email: "alice@example.com", Role: "admin"},
	{ID: 2, Name: "Bob Smith", Email: "bob@example.com", Role: "developer"},
	{ID: 3, Name: "Charlie Brown", Email: "charlie@example.com", Role: "designer"},
	{ID: 4, Name: "Diana Prince", Email: "diana@example.com", Role: "devops"},
}

// enableCORS adds CORS headers to responses
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// handleHealth returns the health status of the service
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "user-service",
	})
}

// handleGetUsers returns all users
func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// handleGetUserByID returns a single user by ID
func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "id query parameter is required"}`, http.StatusBadRequest)
		return
	}

	var id int
	fmt.Sscanf(idStr, "%d", &id)

	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/api/user-health", enableCORS(handleHealth))
	http.HandleFunc("/api/users", enableCORS(handleGetUsers))
	http.HandleFunc("/api/user", enableCORS(handleGetUserByID))

	log.Printf("User Service starting on 0.0.0.1:%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
