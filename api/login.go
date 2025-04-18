package api

import (
	"encoding/json"
	"insights/database"
	"log"
	"net"
	"net/http"
	"time"
)

//Note: The tenant should ideally be identified via the API key / JWT Token, etc
//but the task mentions it comes from the request body
type LoginEvent struct {
	Tenant	  string 	`json:"tenant"`
	User      string 	`json:"user"`
	Origin    string 	`json:"origin"`
	Status    string 	`json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// POST /api/login/new
// Stores details about a login attempt
func NewLogin(w http.ResponseWriter, r *http.Request) {
	// Make sure the request details are valid
	var event LoginEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	// Validate the request details aren't empty
	if event.Tenant == "" || event.User == "" || event.Origin == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	//perhaps other states, but for now only success and failure
	if event.Status != "success" && event.Status != "failure" {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	// Set the timestamp to the current time if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// ensure the ip address valid
	if net.ParseIP(event.Origin) == nil {
		http.Error(w, "Invalid IP address", http.StatusBadRequest)
		return
	}

	// Store the login event in the database
    db, err := database.Connect()
    if err != nil {
        log.Printf("Failed to connect to database: %v", err)
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }

    defer db.Close()
    
	//The requirements want us to make sure the event is idempotent (don't store twice)
    checkSQL := `
        SELECT COUNT(1) FROM login_events
        WHERE tenant = ? AND user = ? AND origin = ? AND status = ? AND timestamp = ?
    `

    var count int
    err = db.QueryRow(checkSQL, event.Tenant, event.User, event.Origin, event.Status, event.Timestamp).Scan(&count)
    if err != nil {
        log.Printf("Failed to check for existing login event: %v", err)
        http.Error(w, "Database query error", http.StatusInternalServerError)
        return
    }

	//Event already exists just say ok and abort out
    if count > 0 {
        w.WriteHeader(http.StatusOK)
        return
    }

    // Insert the login event otherwise
    insertSQL := `INSERT INTO login_events (tenant, user, origin, status, timestamp) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, event.Tenant, event.User, event.Origin, event.Status, event.Timestamp)
    
	//Make sure we all good or let the user know
	if err != nil {
        log.Printf("Failed to store login event: %v", err)
        http.Error(w, "Failed to store login event", http.StatusInternalServerError)
        return
    }

	//Success response
	w.WriteHeader(http.StatusCreated)
}