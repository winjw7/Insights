package main

import (
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

	//testing route
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	//start the server
	http.ListenAndServe(":3000", r)
}