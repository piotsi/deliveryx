package models

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

// Parse the html templates
var templates = template.Must(template.ParseFiles("templates/restaurants.html", "templates/order.html"))

// IndexHandler handles /
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "index")
}

// RestaurantsHandler handles /order
func RestaurantsHandler(w http.ResponseWriter, r *http.Request) {
	rests, err := GetRestaurant()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "restaurants.html", rests) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate(): %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// OrderHandler handles /order/(RestLink)
func OrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	RestLink := vars["RestLink"]

	menu, err := GetMenu(RestLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "order.html", menu) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
