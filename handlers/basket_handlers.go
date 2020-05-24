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
	Items       map[Item]int
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
		http.Redirect(response, request, "/login", http.StatusFound)
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
		basket.Items = make(map[Item]int)
		session.Values["basket"] = basket
	}

	// Retrieve basket
	basket = session.Values["basket"].(*Basket)

	// Check if the restaurant from which items are in the basket is different than the one that user is adding from, if so, then clear the basket
	prevRestLink := basket.RestLink
	if prevRestLink != details.RestLink {
		basket.RestLink = details.RestLink
		basket.Items = make(map[Item]int)
	}

	item := &Item{
		ItemName:  details.ItemName,
		ItemPrice: details.ItemPrice,
		ItemLink:  details.ItemLink,
	}

	// Check if item is in the basket
	// Add item to the basket, if it is inside, increase its amount
	_, isInBasket := basket.Items[*item]
	if !isInBasket {
		basket.Items[*item] = 1
	} else {
		basket.Items[*item]++
	}

	// Update total amount for items in the basket
	CalculateTotalAmount(basket)

	// Save updated basket
	session.Values["basket"] = basket

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

// BasketEmpty empties the basket
func BasketEmpty(response http.ResponseWriter, request *http.Request) {
	// Get a session
	session, err := store.Get(request, "session")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	basket := GetBasket(request)
	basket.Items = make(map[Item]int)
	basket.TotalAmount = "0.00"

	// Save updated basket
	session.Values["basket"] = basket

	err = session.Save(request, response)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(response, request, fmt.Sprintf("/order/%s", basket.RestLink), http.StatusFound)
}

// BasketRemove removes item from the basket or
func BasketRemove(response http.ResponseWriter, request *http.Request) {
	// Get a session
	session, err := store.Get(request, "session")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
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

	item := &Item{
		ItemName:  details.ItemName,
		ItemPrice: details.ItemPrice,
		ItemLink:  details.ItemLink,
	}

	// Retrieve basket
	basket := session.Values["basket"].(*Basket)

	// Remove item or decrease its amount
	if basket.Items[*item] > 1 {
		basket.Items[*item]--
	} else {
		delete(basket.Items, *item)
	}

	// Update total amount for items in the basket
	CalculateTotalAmount(basket)

	// Save updated basket
	session.Values["basket"] = basket

	err = session.Save(request, response)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(response, request, fmt.Sprintf("/order/%s", basket.RestLink), http.StatusFound)
}

// BasketSendOrder sends users order to the restaurant and clears their basket
func BasketSendOrder(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/account", http.StatusFound)
}

// CalculateTotalAmount calculates total amount for prducts in the basket
func CalculateTotalAmount(basket *Basket) {
	total := 0.00
	// itemPriceFloat, err := strconv.ParseFloat(details.ItemPrice, 32)
	for item, amount := range basket.Items {
		itemPriceFloat, _ := strconv.ParseFloat(item.ItemPrice, 32)
		itemTotal := float64(amount) * itemPriceFloat
		total += itemTotal
	}
	basket.TotalAmount = fmt.Sprintf("%.2f", total)
}
