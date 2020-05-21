package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/mitchellh/mapstructure"
)

// Basket holds items added to basket from specified restaurant
type Basket struct {
	RestLink    string
	Items       []Item
	TotalAmount string
	UserName    string
}

// BasketAdd adds items to the basket
func BasketAdd(response http.ResponseWriter, request *http.Request) {
	// Get a session
	session, err := store.Get(request, "session")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect if the user is not logged interface{}
	if !IsAuthenticated(request) {
		http.Redirect(response, request, "/signin", http.StatusFound)
		return
	}

	// Get details from the request
	details := new(Item)
	err = request.ParseForm() // Parse POST form into request.PostForm and request.Form
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(details, request.PostForm) // Decode details from POST form of request to restaurant instance
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	basket := &Basket{}

	// If no basket exists for the session then create one
	if session.Values["basket"] == nil {
		basket.RestLink = details.RestLink
		basket.TotalAmount = "0.00"
		basket.UserName = GetUserName(request)
		session.Values["basket"] = basket
	}

	// Retrieve basket
	basket = session.Values["basket"].(*Basket)
	prevRestLink := basket.RestLink
	if prevRestLink != details.RestLink {
		basket.RestLink = details.RestLink
		basket.Items = nil
		basket.TotalAmount = "0.00"
		log.Printf("Cleared basket for %s", GetUserName(request))
	}

	item := &Item{
		ItemName:  details.ItemName,
		ItemPrice: details.ItemPrice,
		ItemLink:  details.ItemLink,
	}

	// Add item to the basket
	basket.Items = append(basket.Items, *item)
	// Update total amount for items in the basket
	CalculateTotalAmount(basket)

	// Save updated basket
	session.Values["basket"] = basket

	log.Println(session.Values["basket"])

	err = session.Save(request, response)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(response, request, fmt.Sprintf("/order/%s", details.RestLink), http.StatusFound)
	// http.Redirect(response, request, "/basket", http.StatusFound)
}

// GetBasket returns basket contents for given session
func GetBasket(request *http.Request) *Basket {
	// Get a session
	session, err := store.Get(request, "session")
	if err != nil {
		log.Println(err.Error())
	}

	basket := new(Basket)

	mapstructure.Decode(session.Values["basket"], basket)

	return basket
}

// CalculateTotalAmount calculates total amount for prducts in the basket
func CalculateTotalAmount(basket *Basket) {
	total := 0.00
	// itemPriceFloat, err := strconv.ParseFloat(details.ItemPrice, 32)
	for _, item := range basket.Items {
		itemPriceFloat, _ := strconv.ParseFloat(item.ItemPrice, 32)
		total += itemPriceFloat
	}
	basket.TotalAmount = fmt.Sprintf("%.2f", total)
}
