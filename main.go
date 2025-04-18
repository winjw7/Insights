package main

import (
	"insights/api"
	"insights/auth"
	"insights/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Set up the database
	if err := database.Setup(); err != nil {
		log.Fatalf("Failed to set up database: %v", err)
	}

	//Setup the api server
	r := mux.NewRouter()

	//Ensure the request is from a valid tennant
	r.Use(auth.Verify)

	//setup the api routes
	r.HandleFunc("/api/login/new", api.NewLogin).Methods("POST")
	r.HandleFunc("/api/login/suspicious", api.GetSuspiciousLogins).Methods("GET")
	
	//start the server
	http.ListenAndServe(":3000", r)
}