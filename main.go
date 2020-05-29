package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"os"
	"site/handlers"

	"github.com/gorilla/mux"
)

func main() {
<<<<<<< HEAD
=======
	// Parse a session key from the flag
	sessionKey := flag.String("sessionkey", "", "Provide a session key")
	flag.Parse()
	handlers.SessionKey = *sessionKey

	if *sessionKey == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

>>>>>>> 5973ac5fd46483de6b52bd96d50aafe366a66f9b
	// Gorilla mux router, StrictSlash() adds trailing slash to the end of the url path
	// Note: I don't know if it is caused by this func but you need to add trailing slash to action attribute link e.g. action="/page/"
	router := mux.NewRouter().StrictSlash(true)

	// Database initialization
	handlers.InitDB("root:mysqlPSWD213@(127.0.0.1:3306)/root")

	// Register basket for the session
	gob.Register(&handlers.Basket{})

	// Handlers, Methods() accepts only matched methods
	router.HandleFunc("/signmein/", handlers.SignIn).Methods("POST")
	router.HandleFunc("/signmeup/", handlers.SignUp).Methods("POST")
	router.HandleFunc("/accountedit/", handlers.AccountEdit).Methods("POST")
	router.HandleFunc("/restaurantedit/", handlers.RestaurantEdit).Methods("POST")
	router.HandleFunc("/logoedit/", handlers.LogoEdit).Methods("POST")
	router.HandleFunc("/itemedit/", handlers.ItemEdit).Methods("POST")
	router.HandleFunc("/itemremove/", handlers.ItemRemove).Methods("POST")
	router.HandleFunc("/itemadd/", handlers.ItemAdd).Methods("POST")
	router.HandleFunc("/basketadd/", handlers.BasketAdd).Methods("POST")
	router.HandleFunc("/basketremove/", handlers.BasketRemove).Methods("POST")
	router.HandleFunc("/ordercomplete/", handlers.OrderComplete).Methods("POST")
	router.HandleFunc("/basketempty/", handlers.BasketEmpty)
	router.HandleFunc("/sendorder/", handlers.BasketSendOrder)
	router.HandleFunc("/signmeout/", handlers.SignMeOut)
	router.HandleFunc("/", handlers.IndexPageHandler)
	router.HandleFunc("/login/", handlers.LoginPageHandler)
	router.HandleFunc("/order/", handlers.RestaurantsPageHandler)
	router.HandleFunc("/order/{RestLink}/", handlers.OrderPageHandler)
	router.HandleFunc("/account/", handlers.AccountPageHandler)
	router.HandleFunc("/orders/", handlers.OrdersPageHandler)

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
