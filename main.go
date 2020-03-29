package main

import (
	"fmt"
	"log"
	"net/http"
	"site/models"
	"text/template"

	"github.com/gorilla/mux"
)

// Parse the html templates
var templates = template.Must(template.ParseFiles("templates/restaurants.html", "templates/order.html"))

func main() {
	// Gorilla mux router
	r := mux.NewRouter().StrictSlash(true) // Add trailing slash to the end of path

	// Database initialization
	models.InitDB("root:mysqlPSWD213@(127.0.0.1:3306)/root")

	// Handlers
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/order/", restaurantsHandler).Methods("GET") // Narrow handling by GET
	r.HandleFunc("/order/{RestLink}/", orderHandler).Methods("GET")

	// Static handlers
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images/")))) // Handle static files in images folder
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))          // Handle static files in css folder

	// Return errors on TCP network
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", r)
	}
}

// restaurantsHandler handles /order
func restaurantsHandler(w http.ResponseWriter, r *http.Request) {
	rests, err := models.GetRestaurant()
	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(500), 500) // Internal server error
		return
	}
	err = templates.ExecuteTemplate(w, "restaurants.html", rests) // Execute parsed template
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// orderHandler handles /order/(RestLink)
func orderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	RestLink := vars["RestLink"]
	menu, err := models.GetMenu(RestLink)
	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(500), 500) // Internal server error
		return
	}
	err = templates.ExecuteTemplate(w, "order.html", menu) // Execute parsed template
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "index")
}
