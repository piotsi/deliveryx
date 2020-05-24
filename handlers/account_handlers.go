package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

// AccountEdit edits account details
func AccountEdit(response http.ResponseWriter, request *http.Request) {
	// Get credentials from the request
	credentials := new(Credentials)
	err := request.ParseForm() // Parse POST form into request.PostForm and request.Form
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(credentials, request.PostForm) // Decode credentials from POST form of request to credentials instance
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encrypt password
	hashedUserPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.UserPassword), 8)

	// Insert credentials into the database
	query := fmt.Sprintf("UPDATE users SET userPassword='%s' WHERE userName='%s'", string(hashedUserPassword), credentials.UserName)
	_, err = db.Query(query)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(response, request, "/account", http.StatusFound)
}

// RestaurantEdit edits restaurant details
func RestaurantEdit(response http.ResponseWriter, request *http.Request) {
	owner := GetUserName(request)

	// Get credentials from the request
	details := new(Restaurant)
	err := request.ParseForm() // Parse POST form into request.PostForm and request.Form
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(details, request.PostForm) // Decode details from POST form of request to restaurant instance
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	// Change details in the database
	query := fmt.Sprintf("UPDATE restaurants SET RestName='%s', RestAddress='%s', RestLink='%s' WHERE RestOwner='%s'", details.RestName, details.RestAddress, details.RestLink, owner)
	_, err = db.Query(query)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(response, request, "/account", http.StatusFound)
}

// LogoEdit changes restaurants' logo
func LogoEdit(response http.ResponseWriter, request *http.Request) {
	// Parse file from request
	err := request.ParseMultipartForm(1)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	imageLink := request.FormValue("imageLink")

	file, handler, err := request.FormFile("imageFile")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create temporary file
	tempFile, err := ioutil.TempFile("images/restaurants/", "*.png")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Read uploaded file into byte arrat
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write read byte array into temporary file
	tempFile.Write(fileBytes)

	// Rename tempfile to old logo name
	// If the file exists it will overwrite old one, if it doesn't exist it will only change name of the new one
	newFilePathName := fmt.Sprintf("images/restaurants/%s.png", imageLink)
	err = os.Rename(tempFile.Name(), newFilePathName)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("%s file uploaded", handler.Filename)
	http.Redirect(response, request, "/account", http.StatusFound)
}

// ItemEdit handles changing, adding and removing items
func ItemEdit(response http.ResponseWriter, request *http.Request) {
	// Get details from the request
	details := new(Item)
	err := request.ParseForm() // Parse POST form into request.PostForm and request.Form
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(details, request.PostForm) // Decode details from POST form of request to restaurant instance
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	restLink := GetRestLink(request)

	// Change details in the database
	query := fmt.Sprintf("UPDATE items SET ItemName='%s', ItemPrice='%s', ItemDescription='%s' WHERE ItemLink='%s' AND RestLink='%s'", details.ItemName, details.ItemPrice, details.ItemDescription, details.ItemLink, restLink)
	_, err = db.Query(query)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(response, request, "/account", http.StatusFound)
}

// ItemRemove removes item specified in the request form
func ItemRemove(response http.ResponseWriter, request *http.Request) {
	// Get details from the request
	details := new(Item)
	err := request.ParseForm() // Parse POST form into request.PostForm and request.Form
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(details, request.PostForm) // Decode details from POST form of request to restaurant instance
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	// Change details in the database
	query := fmt.Sprintf("DELETE FROM items WHERE ItemLink='%s'", details.ItemLink)
	_, err = db.Query(query)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(response, request, "/account", http.StatusFound)
}

// ItemAdd adds new item specified in the request form
func ItemAdd(response http.ResponseWriter, request *http.Request) {
	// Get details from the request
	details := new(Item)
	err := request.ParseForm() // Parse POST form into request.PostForm and request.Form
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(details, request.PostForm) // Decode details from POST form of request to restaurant instance
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	itemLink := GenerateItemLink(request, details.ItemName)
	restLink := GetRestLink(request)

	// Change details in the database
	query := fmt.Sprintf("INSERT INTO items (ItemName, ItemPrice, ItemDescription, ItemLink, RestLink) VALUES ('%s', '%s', '%s', '%s', '%s')", details.ItemName, details.ItemPrice, details.ItemDescription, itemLink, restLink)
	_, err = db.Query(query)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(response, request, "/account", http.StatusFound)
}

// GetRestaurantDetails obtains restaurant details of which logged user is the owner of
func GetRestaurantDetails(request *http.Request) *Restaurant {
	RestDetails := new(Restaurant)
	owner := GetUserName(request)

	query := fmt.Sprintf("SELECT RestName, RestAddress, RestLink FROM restaurants WHERE RestOwner='%s'", owner)
	row := db.QueryRow(query)

	row.Scan(&RestDetails.RestName, &RestDetails.RestAddress, &RestDetails.RestLink)

	return RestDetails
}

// GetRestLink returns RestLink for which logged user is the owner
func GetRestLink(request *http.Request) string {
	var RestLink string

	owner := GetUserName(request)

	query := fmt.Sprintf("SELECT RestLink FROM restaurants WHERE RestOwner='%s'", owner)
	row := db.QueryRow(query)

	err := row.Scan(&RestLink)
	if err != nil {
		return ""
	}

	return RestLink
}

// GenerateItemLink generates link for the new item
func GenerateItemLink(request *http.Request, name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "-", -1)
	name = GetRestLink(request) + "-" + name
	return name
}
