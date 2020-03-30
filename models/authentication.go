package models

import (
    "fmt"
	"encoding/json"
	"net/http"
    "golang.org/x/crypto/bcrypt"
    "database/sql"
)

// Credentials stores authentication data
type Credentials struct {
	UserName       string `json:"userName" db:"userName"`
	UserPassword   string `json:"userPassword" db:"userPassword"`
}

// SignIn handles logging in
func SignIn(w http.ResponseWriter, r *http.Request) {
    // Get credentials from the request
    credentials := new(Credentials)
    err := json.NewDecoder(r.Body).Decode(credentials)  // Decode request body to credentials instance
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    // Get password from the database
    query := fmt.Sprintf("SELECT userPassword FROM users WHERE userName='%s'", credentials.UserName)
    result := db.QueryRow(query)

    storedCredentials := new(Credentials)
    err = result.Scan(&storedCredentials.UserPassword)
	if err != nil {
        // If UserName doesn't exist
		if err == sql.ErrNoRows {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    // Compare hashedUserPassword from databse and from the request
    err = bcrypt.CompareHashAndPassword([]byte(storedCredentials.UserPassword), []byte(credentials.UserPassword))
    if err != nil {
        // Return error if passwords don't match
        http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}

// SignUp handles creating new account
func SignUp(w http.ResponseWriter, r *http.Request) {
    // Get credentials from the request
	credentials := new(Credentials)
	err := json.NewDecoder(r.Body).Decode(credentials)  // Decode request body to credentials instance
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    hashedUserPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.UserPassword), 8)

    // Insert credentials into the database
    query := fmt.Sprintf("INSERT INTO users(userName, userPassword) VALUES ('%s', '%s')", credentials.UserName, string(hashedUserPassword))
    _, err = db.Query(query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
