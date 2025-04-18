package auth

import (
	"net/http"
)

const apiKeyHeader string = "X-API-Key";

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

		//Move on to other middleware etc
		next.ServeHTTP(w, r)
	})
}