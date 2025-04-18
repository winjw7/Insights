package api

import (
	"encoding/json"
	"insights/auth"
	"insights/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

// What the user gets back from the API
type SuspiciousLoginResponse struct {
    Origin      string `json:"origin"`
    FailCount   int    `json:"failCount"`
}

const defaultThreshold = 5
const defaultMinutes = 3
const defaultPage = 1
const defaultLimit = 10
const maxLimit = 100
const minLimit = 1

// Note: The login request tenant to be submitted via the request body (insecurely),
// however objective 4 wants me to ensure each tenant can only see their own data
// so for the GET we will actually use the "API Key" to try and
// fufill both objectives.

// GET /api/login/suspicious
func GetSuspiciousLogins(w http.ResponseWriter, r *http.Request) {

    //grab tennant id securely from middleware
    tenantID := auth.GetTenantID(r)

    //This should never happen because we passed middleware, but safety :)
    if tenantID == "" {
        http.Error(w, "Tenant ID is missing", http.StatusUnauthorized)
        return
    }

    // Get query parameters with default values
    threshold, err := getIntParam(r, "threshold", defaultThreshold, 1, 0)
    if err != nil {
        http.Error(w, "Invalid threshold value", http.StatusBadRequest)
        return
    }

    minutes, err := getIntParam(r, "minutes", defaultMinutes, 1, 0)
    if err != nil {
        http.Error(w, "Invalid minutes value", http.StatusBadRequest)
        return
    }

    page, err := getIntParam(r, "page", defaultPage, 1, 0)
    if err != nil {
        http.Error(w, "Invalid page value", http.StatusBadRequest)
        return
    }

    limit, err := getIntParam(r, "limit", defaultLimit, minLimit, maxLimit)
    if err != nil {
        http.Error(w, "Invalid limit value (must be between 1 and 100)", http.StatusBadRequest)
        return
    }

    order := "DESC" // Default order
    orderParam := r.URL.Query().Get("order")

    if orderParam == "asc" {
        order = "ASC"
    }

    // Connect to database
    db, err := database.Connect()
    if err != nil {
        log.Printf("Failed to connect to database: %v", err)
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // Calculate the time window
    timeWindow := time.Now().Add(time.Duration(-minutes) * time.Minute)

    // query stuff
    query := `
        SELECT origin, COUNT(*) as fail_count
        FROM login_events
        WHERE tenant = ? AND status = 'failure' AND timestamp > ?
        GROUP BY origin
        HAVING COUNT(*) >= ?
        ORDER BY fail_count ` + order + `
        LIMIT ? OFFSET ?
    `

    offset := (page - 1) * limit
    rows, err := db.Query(query, tenantID, timeWindow, threshold, limit, offset)

    if err != nil {
        log.Printf("Failed to execute query: %v", err)
        http.Error(w, "Database query error", http.StatusInternalServerError)
        return
    }

    defer rows.Close()

    var results [] SuspiciousLoginResponse

    //grab rows and add to results
    for rows.Next() {
        var result SuspiciousLoginResponse

        if err := rows.Scan(&result.Origin, &result.FailCount); err != nil {
            log.Printf("Failed to scan row: %v", err)
            http.Error(w, "Database scan error", http.StatusInternalServerError)
            return
        }

        results = append(results, result)
    }

    //ensure rows didnt have any problems
    if err := rows.Err(); err != nil {
        log.Printf("Row error: %v", err)
        http.Error(w, "Database row error", http.StatusInternalServerError)
        return
    }

    //instead of null send back empty array
    if(len(results) == 0) {
        results = []SuspiciousLoginResponse{}
    }

    //Give results back as json
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    if err := json.NewEncoder(w).Encode(results); err != nil {
        log.Printf("Failed to encode response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}

// Helper function to parse int from query
func getIntParam(r *http.Request, name string, defaultValue, minValue, maxValue int) (int, error) {
    
    strVal := r.URL.Query().Get(name)
    
    //no value given, so go with default
    if strVal == "" {
        return defaultValue, nil
    }
    
    // attempt to parse the value
    val, err := strconv.Atoi(strVal)
    if err != nil {
        return 0, err
    }
    
    // check min (if set)
    if minValue != 0 && val < minValue {
        return 0, err
    }
    
    //check max (if set)
    if maxValue != 0 && val > maxValue {
        return 0, err
    }
    
    return val, nil
}