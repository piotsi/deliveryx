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
    fmt.Println(credentials.UserName)
    query := fmt.Sprintf("UPDATE users SET userPassword='%s' WHERE userName='%s'", string(hashedUserPassword), credentials.UserName)
    _, err = db.Query(query)
    if err != nil {
        http.Error(response, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(response, request, "/", http.StatusFound)
}
