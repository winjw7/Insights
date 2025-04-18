package auth

import (
	"context"
	"net/http"
)

const apiKeyHeader string = "X-API-Key";

type contextKey string
const tennantIDHeader contextKey = "Tennant-ID"

// Ensure a request is from an authorized tennant 
func Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		//Get the API key from the request header
		apiKey := r.Header.Get(apiKeyHeader)
		
		//Ensure api key actually exists
		if apiKey == "" {
			http.Error(w, "API key is missing", http.StatusUnauthorized)
			return
		}

		//Is the API key valid now
		if !authenticate(apiKey) {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		//Add the tennant id we "grabbed" from the API key to the request context
		ctx := context.WithValue(r.Context(), tennantIDHeader, apiKey)
		r = r.WithContext(ctx)

		//Move on to other middleware etc
		next.ServeHTTP(w, r)
	})
}

// GetTenantID extracts the tenant ID from the request context
func GetTenantID(r *http.Request) string {
    if tenantID, ok := r.Context().Value(tennantIDHeader).(string); ok {
        return tenantID
    }
    return ""
}