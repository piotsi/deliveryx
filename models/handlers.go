package models

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

// Parse the html templates
var templates = template.Must(template.ParseFiles(
	"templates/restaurants.html",
	"templates/order.html",
	"templates/signin.html",
	"templates/signup.html"))

// IndexHandler handles /
func IndexHandler(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/order", http.StatusFound)
}

// SigninHandler handles /signin page
func SigninHandler(response http.ResponseWriter, request *http.Request) {
	err := templates.ExecuteTemplate(response, "signin.html", "") // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate(): %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

// SignupHandler handles /signup page
func SignupHandler(response http.ResponseWriter, request *http.Request) {
	err := templates.ExecuteTemplate(response, "signup.html", "") // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate(): %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

// RestaurantsHandler handles /order page
func RestaurantsHandler(response http.ResponseWriter, request *http.Request) {
	rests, err := GetRestaurant()
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(response, "restaurants.html", rests) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate(): %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

// OrderHandler handles /order/(RestLink) page
func OrderHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	RestLink := vars["RestLink"]

	menu, err := GetMenu(RestLink)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(response, "order.html", menu) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate: %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}
