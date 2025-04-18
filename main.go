package main

import (
	"insights/api"
	"insights/auth"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	//Setup the api server
	r := mux.NewRouter()

	//Ensure the request is from a valid tennant
	r.Use(auth.Verify)

	//setup the api routes
	r.HandleFunc("/api/login/new", api.NewLogin).Methods("POST")

	//start the server
	http.ListenAndServe(":3000", r)
}