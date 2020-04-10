package handlers

import (
    "os"
    "fmt"
	"net/http"
    "golang.org/x/crypto/bcrypt"
    "database/sql"
    "log"
    "github.com/gorilla/sessions"
    "github.com/gorilla/schema"
    // "github.com/gorilla/securecookie"
)

// Credentials stores authentication data
type Credentials struct {
	UserName     string `db:"userName"`
	UserPassword string `db:"userPassword"`
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

// SignIn handles /signmein
func SignIn(response http.ResponseWriter, request *http.Request) {
    session, _ := store.Get(request, "cookie-name")

    // Get credentials from the request
    credentials := new(Credentials)
    err := request.ParseForm() // Parse POST form into request.PostForm and request.Form
    if err != nil {
        log.Fatal(err.Error())
    }
    decoder := schema.NewDecoder()
    err = decoder.Decode(credentials, request.PostForm) // Decode credentials from POST form of request to credentials instance
    if err != nil {
        log.Fatal(err.Error())
    }

    session.Values["authenticated"] = true
    session.Values["username"] = credentials.UserName

    // Get password from the database
    query := fmt.Sprintf("SELECT userPassword FROM users WHERE userName='%s'", credentials.UserName)
    result := db.QueryRow(query)

    storedCredentials := &Credentials{}
    err = result.Scan(&storedCredentials.UserPassword)
	if err != nil {
        // If UserName doesn't exist
		if err == sql.ErrNoRows {
            http.Error(response, "User doesn't exist!", http.StatusInternalServerError)
            return
		}
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

    // Compare hashedUserPassword from databse and from the request
    err = bcrypt.CompareHashAndPassword([]byte(storedCredentials.UserPassword), []byte(credentials.UserPassword))
    if err != nil {
        // Return error if passwords don't match
        http.Error(response, "Wrong password!", http.StatusUnauthorized)
        return
	}
    session.Save(request, response)
    log.Printf("%s logged in!", credentials.UserName)
    http.Redirect(response, request, "/", http.StatusFound)
}

// SignUp handles /signmeup
func SignUp(response http.ResponseWriter, request *http.Request) {
    // session, _ := store.Get(request, "cookie-name")

    // Get credentials from the request
    credentials := new(Credentials)
    err := request.ParseForm() // Parse POST form into request.PostForm and request.Form
    if err != nil {
        log.Fatal(err.Error())
    }
    decoder := schema.NewDecoder()
    err = decoder.Decode(credentials, request.PostForm) // Decode credentials from POST form of request to credentials instance
    if err != nil {
        log.Fatal(err.Error())
    }

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

// SignMeOut handles /signmeout
func SignMeOut(response http.ResponseWriter, request *http.Request) {
    session, _ := store.Get(request, "cookie-name")

    session.Values["authenticated"] = false
    log.Printf("%s logged out!", session.Values["username"].(string))

    session.Save(request, response)
    http.Redirect(response, request, "/", http.StatusFound)
}

// GetUserName returns authenticated user's UserName
func GetUserName(request *http.Request) (UserName string) {
    session, _ := store.Get(request, "cookie-name")
    if session.Values["authenticated"] == true {
        return session.Values["username"].(string)
    }
    return
}
