package models

import (
    "fmt"
	"net/http"
    "golang.org/x/crypto/bcrypt"
    "database/sql"
    "log"
)

// Credentials stores authentication data
type Credentials struct {
	UserName     string `db:"userName"`
	UserPassword string `db:"userPassword"`
}

// SignIn handles logging in
func SignIn(response http.ResponseWriter, request *http.Request) {
    // Get credentials from the request
    credentials := &Credentials{UserName: request.FormValue("userName"), UserPassword: request.FormValue("userPassword")}

    // Get password from the database
    query := fmt.Sprintf("SELECT userPassword FROM users WHERE userName='%s'", credentials.UserName)
    result := db.QueryRow(query)

    storedCredentials := &Credentials{}
    err := result.Scan(&storedCredentials.UserPassword)
	if err != nil {
        // If UserName doesn't exist
		if err == sql.ErrNoRows {
            fmt.Fprintf(response, "User doesn't exist!")
            return
		}
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

    // Compare hashedUserPassword from databse and from the request
    err = bcrypt.CompareHashAndPassword([]byte(storedCredentials.UserPassword), []byte(credentials.UserPassword))
    if err != nil {
        // Return error if passwords don't match
        // http.Error(response, err.Error(), http.StatusUnauthorized)
        fmt.Fprintf(response, "Wrong password!")
	}
    log.Printf("%s logged in!", credentials.UserName)
    http.Redirect(response, request, "/", http.StatusFound)
}

// SignUp handles creating new account
func SignUp(response http.ResponseWriter, request *http.Request) {
    // Get credentials from the request
    credentials := &Credentials{UserName: request.FormValue("userName"), UserPassword: request.FormValue("userPassword")}

    // Encrypt password
    hashedUserPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.UserPassword), 8)

    // Insert credentials into the database
    query := fmt.Sprintf("INSERT INTO users(userName, userPassword) VALUES ('%s', '%s')", credentials.UserName, string(hashedUserPassword))
    _, err = db.Query(query)
    if err != nil {
        http.Error(response, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(response, request, "/", http.StatusFound)
}
