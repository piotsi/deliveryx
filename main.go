package main

import (
	"log"
	"net/http"
	"site/models"

	"github.com/gorilla/mux"
)

func main() {
	// Gorilla mux router
	r := mux.NewRouter().StrictSlash(true) // Add trailing slash to the end of path

	// Database initialization
	models.InitDB("root:mysqlPSWD213@(127.0.0.1:3306)/root")

	// Handlers
	r.HandleFunc("/", models.IndexHandler)
	r.HandleFunc("/signin/", models.SignIn).Methods("POST") // Accept only POST request
	r.HandleFunc("/signup/", models.SignUp).Methods("POST")
	r.HandleFunc("/order/", models.RestaurantsHandler).Methods("GET") // Accept only GET request
	r.HandleFunc("/order/{RestLink}/", models.OrderHandler).Methods("GET")

	// Static handlers
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images/")))) // Handle static files in images folder
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))          // Handle static files in css folder

	// Return errors on TCP network
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("http.ListenAndServe: ", r)
	}
}
