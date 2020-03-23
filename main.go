package main

import (
	"fmt"
	"log"
	"net/http"
	"site/models"
	"text/template"

	"github.com/gorilla/mux"
)

type data struct {
	Title string
}

var templates = template.Must(template.ParseFiles("templates/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, d *data) {
	err := templates.ExecuteTemplate(w, tmpl+".html", d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handle favicon
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "images/favicon.ico")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	d := data{
		Title: "index",
	}
	if r.Method == "POST" {
		http.Redirect(w, r, "order", http.StatusSeeOther)
	}
	renderTemplate(w, "index", &d)
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "order")

}

func main() {
	// Gorilla mux router
	r := mux.NewRouter()

	// Database initialization
	models.InitDB("root:mysqlPSWD213@(127.0.0.1:3306)/root?parseTime=true")

	// Handlers
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.HandleFunc("/order", orderHandler)
	r.HandleFunc("/", indexHandler)

	// Return errors on TCP network
	lasErr := http.ListenAndServe(":8080", r)
	if lasErr != nil {
		log.Fatal("ListenAndServe: ", lasErr)
	}
}
