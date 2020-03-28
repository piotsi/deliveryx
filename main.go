package main

import (
	"log"
	"net/http"
	"site/models"
	"text/template"

	"github.com/gorilla/mux"
)

// Parse the html templates
var templates = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	// Gorilla mux router
	r := mux.NewRouter()

	// Database initialization
	models.InitDB("root:mysqlPSWD213@(127.0.0.1:3306)/root")

	// Handlers
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/order", restaurantsHandler)

	r.Handle("/css/{style}", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))                                              // Handle static files
	r.Handle("/images/{image}", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))                                     // Handle logo and favicon
	r.Handle("/images/restaurants/{image}", http.StripPrefix("/images/restaurants/", http.FileServer(http.Dir("images/restaurants/")))) // Handle restaurant images                                                                       // Handle favicon

	// Return errors on TCP network
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", r)
	}
}

func restaurantsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405) // Method not allowed
		return
	}
	rests, err := models.GetRestaurant()
	if err != nil {
		http.Error(w, http.StatusText(500), 500) // Internal server error
		return
	}
	// for _, rest := range rests {
	// 	fmt.Fprintf(w, "%s, %s, %s\n", rest.RestName, rest.RestAddress, rest.RestPictureLocation)
	// }
	err = templates.ExecuteTemplate(w, "index.html", rests) // Execute parsed template
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// d := data{
	// 	Title: "index",
	// }
	// if r.Method == "POST" {
	// 	http.Redirect(w, r, "order", http.StatusSeeOther)
	// }
	// renderTemplate(w, "index", &d)
}
