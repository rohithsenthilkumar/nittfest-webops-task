package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Order represents an order entity
type Order struct {
	ID       int    `json:"id"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name,omitempty"`
	Status   string `json:"status"`
}

// User represents user data fetched from user-service
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Hardcoded orders (no database)
var orders = []Order{
	{ID: 101, Item: "Laptop", Quantity: 1, UserID: 1, Status: "shipped"},
	{ID: 102, Item: "Keyboard", Quantity: 2, UserID: 2, Status: "processing"},
	{ID: 103, Item: "Monitor", Quantity: 1, UserID: 3, Status: "delivered"},
	{ID: 104, Item: "Mouse", Quantity: 3, UserID: 1, Status: "processing"},
	{ID: 105, Item: "Headphones", Quantity: 1, UserID: 4, Status: "shipped"},
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

// getUserServiceURL returns the user-service base URL
func getUserServiceURL() string {
	url := os.Getenv("USER_SERVICE_URL")
	if url == "" {
		url = "http://localhost:8081"
	}
	return url
}

// fetchUserName calls user-service to get a user's name by ID
func fetchUserName(userID int) string {
	url := fmt.Sprintf("%s/api/user?id=%d", getUserServiceURL(), userID)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("  Failed to reach user-service: %v", err)
		return "Unknown"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Unknown"
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return "Unknown"
	}

	return user.Name
}

// handleHealth returns the health status of the service
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "order-service",
	})
}

// handleGetOrders returns all orders enriched with user names
func handleGetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Enrich orders with user names by calling user-service
	enrichedOrders := make([]Order, len(orders))
	copy(enrichedOrders, orders)

	for i := range enrichedOrders {
		enrichedOrders[i].UserName = fetchUserName(enrichedOrders[i].UserID)
	}

	json.NewEncoder(w).Encode(enrichedOrders)
}

// handleGetOrderByID returns a single order by ID, enriched with user name
func handleGetOrderByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "id query parameter is required"}`, http.StatusBadRequest)
		return
	}

	var id int
	fmt.Sscanf(idStr, "%d", &id)

	for _, order := range orders {
		if order.ID == id {
			order.UserName = fetchUserName(order.UserID)
			json.NewEncoder(w).Encode(order)
			return
		}
	}

	http.Error(w, `{"error": "order not found"}`, http.StatusNotFound)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Order Service starting on 0.0.0.1:%s", port)
	log.Printf("User Service URL: %s", getUserServiceURL())

	http.HandleFunc("/api/order-health", enableCORS(handleHealth))
	http.HandleFunc("/api/orders", enableCORS(handleGetOrders))
	http.HandleFunc("/api/order", enableCORS(handleGetOrderByID))

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
