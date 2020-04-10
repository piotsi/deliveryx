package handlers

import (
	"fmt"
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

// IndexPageHandler handles /
func IndexPageHandler(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/order", http.StatusFound)
}

// SigninPageHandler handles /signin page
func SigninPageHandler(response http.ResponseWriter, request *http.Request) {
	err := templates.ExecuteTemplate(response, "signin.html", "") // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate(): %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

// SignupPageHandler handles /signup page
func SignupPageHandler(response http.ResponseWriter, request *http.Request) {
	err := templates.ExecuteTemplate(response, "signup.html", "") // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate(): %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

// RestaurantsPageHandler handles /order page
func RestaurantsPageHandler(response http.ResponseWriter, request *http.Request) {
	rests, err := GetRestaurants()
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(response, "restaurants.html", map[string]interface{}{"Username": GetUserName(request), "Rest": rests}) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate(): %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

// OrderPageHandler handles /order/(RestLink) page
func OrderPageHandler(response http.ResponseWriter, request *http.Request) {
	routeVars := mux.Vars(request)
	restLink := routeVars["RestLink"]

	menu, err := GetMenu(restLink)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(response, "order.html", map[string]interface{}{"Username": GetUserName(request), "Menu": menu}) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate: %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

// AccountPageHandler handles /account page
func AccountPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "User logged: %s", GetUserName(request))
}
