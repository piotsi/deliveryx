package main

import (
	"log"
	"net/http"
	"site/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Gorilla mux router, StrictSlash() adds trailing slash to the end of path
	router := mux.NewRouter().StrictSlash(true)

	// Database initialization
	handlers.InitDB("root:mysqlPSWD213@(127.0.0.1:3306)/root")

	// Handlers, Methods() accepts only matched methods
	router.HandleFunc("/signmein/", handlers.SignIn).Methods("POST")
	router.HandleFunc("/signmeup/", handlers.SignUp).Methods("POST")
	router.HandleFunc("/accountedit/", handlers.AccountEdit).Methods("POST")
	router.HandleFunc("/restaurantedit/", handlers.RestaurantEdit).Methods("POST")
	router.HandleFunc("/signmeout/", handlers.SignMeOut)
	router.HandleFunc("/", handlers.IndexPageHandler)
	router.HandleFunc("/signin/", handlers.SigninPageHandler)
	router.HandleFunc("/signup/", handlers.SignupPageHandler)
	router.HandleFunc("/order/", handlers.RestaurantsPageHandler)
	router.HandleFunc("/order/{RestLink}/", handlers.OrderPageHandler)
	router.HandleFunc("/account/", handlers.AccountPageHandler)

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
