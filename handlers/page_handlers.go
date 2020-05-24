package handlers

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
	"templates/login.html",
	"templates/account.html",
	"templates/orders.html"))

// IndexPageHandler handles /
func IndexPageHandler(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/order", http.StatusFound)
}

// LoginPageHandler handles /login page
func LoginPageHandler(response http.ResponseWriter, request *http.Request) {
	err := templates.ExecuteTemplate(response, "login.html", map[string]interface{}{"Username": GetUserName(request)}) // Execute parsed template
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

	basket := GetBasket(request)

	err = templates.ExecuteTemplate(response, "restaurants.html", map[string]interface{}{"Username": GetUserName(request), "Restaurant": GetRestaurantDetails(request), "Rest": rests, "Basket": basket}) // Execute parsed template
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

	basket := GetBasket(request)

	err = templates.ExecuteTemplate(response, "order.html", map[string]interface{}{"Username": GetUserName(request), "Restaurant": GetRestaurantDetails(request), "Menu": menu, "Basket": basket, "RestName": GetRestName(request)}) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate: %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

// AccountPageHandler handles /account page
func AccountPageHandler(response http.ResponseWriter, request *http.Request) {
	if !IsAuthenticated(request) {
		http.Error(response, "Forbidden", http.StatusForbidden)
		return
	}

	RestLink := GetRestLink(request)
	item, err := GetMenu(RestLink)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	basket := GetBasket(request)

	err = templates.ExecuteTemplate(response, "account.html", map[string]interface{}{"Username": GetUserName(request), "Restaurant": GetRestaurantDetails(request), "Item": item, "Basket": basket}) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate: %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

// OrdersPageHandler handles /orders page
func OrdersPageHandler(response http.ResponseWriter, request *http.Request) {
	if !IsAuthenticated(request) {
		http.Error(response, "Forbidden", http.StatusForbidden)
		return
	}

	err := templates.ExecuteTemplate(response, "orders.html", map[string]interface{}{"Username": GetUserName(request), "Restaurant": GetRestaurantDetails(request)}) // Execute parsed template
	if err != nil {
		log.Fatalf("templates.ExecuteTemplate: %s", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}
