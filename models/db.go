package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"log"
)

var db *sql.DB

// InitDB initializes MySQL database, takes source name and returns nothing.
func InitDB(source string) {
	db, openErr := sql.Open("mysql", source)
	if openErr != nil {
		log.Fatal(openErr)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
}
