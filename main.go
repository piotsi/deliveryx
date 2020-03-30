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
	r.HandleFunc("/signin/", models.SignIn).Methods("POST") // Accept only POST request
	r.HandleFunc("/signup/", models.SignUp)
	r.HandleFunc("/order/", restaurantsHandler).Methods("GET") // Accept only GET request
	r.HandleFunc("/order/{RestLink}/", orderHandler).Methods("GET")

	// Static handlers
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images/")))) // Handle static files in images folder
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))          // Handle static files in css folder

	// Return errors on TCP network
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("http.ListenAndServe: ", r)
	}
}

// indexHandler handles /
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "index")
}

// restaurantsHandler handles /order
func restaurantsHandler(w http.ResponseWriter, r *http.Request) {
	rests, err := models.GetRestaurant()
	if err != nil {
		log.Fatalf("models.GetRestaurant(): %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "restaurants.html", rests) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate(): %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// orderHandler handles /order/(RestLink)
func orderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	RestLink := vars["RestLink"]

	menu, err := models.GetMenu(RestLink)
	if err != nil {
		log.Fatalf("models.GetMenu(): %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "order.html", menu) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
