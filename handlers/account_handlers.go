package handlers

import (
    "fmt"
    "net/http"
    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/schema"
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
    http.Redirect(response, request, "/", http.StatusFound)
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
    http.Redirect(response, request, "/", http.StatusFound)
}

// GetRestaurantDetails obtains restaurant details of which logged user is the owner of
func GetRestaurantDetails(request *http.Request) (*Restaurant) {
    RestDetails := new(Restaurant)
    owner := GetUserName(request)

    query := fmt.Sprintf("SELECT RestName, RestAddress, RestLink FROM restaurants WHERE RestOwner='%s'", owner)
    row := db.QueryRow(query)

    row.Scan(&RestDetails.RestName, &RestDetails.RestAddress, &RestDetails.RestLink)

    return RestDetails
}
