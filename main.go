package main

import (
	"log"
	"net/http"
	"site/models"

	"github.com/gorilla/mux"
)

func main() {
	// Gorilla mux router, StrictSlash() adds trailing slash to the end of path
	router := mux.NewRouter().StrictSlash(true)

	// Database initialization
	models.InitDB("root:mysqlPSWD213@(127.0.0.1:3306)/root")

	// Handlers, Methods() accepts only matched methods
	router.HandleFunc("/signmein/", models.SignIn).Methods("POST")
	router.HandleFunc("/signmeup/", models.SignUp).Methods("POST")
	router.HandleFunc("/", models.IndexHandler)
	router.HandleFunc("/signin/", models.SigninHandler)
	router.HandleFunc("/signup/", models.SignupHandler)
	router.HandleFunc("/order/", models.RestaurantsHandler)
	router.HandleFunc("/order/{RestLink}/", models.OrderHandler)

	// Static handlers
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images/")))) // Handle static files in images folder
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))          // Handle static files in css folder

	http.Handle("/", router)

	// Return errors on TCP network
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("http.ListenAndServe: ", router)
	}
}
