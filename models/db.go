package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var db *sql.DB

// InitDB initializes MySQL database, takes source name and returns nothing.
func InitDB(source string) {
	var err error // Short declaration needed here! Assignment in sql.Open() results in "runtime error: invalid memory address or nil pointer dereference"
	db, err = sql.Open("mysql", source)
	if err != nil {
		log.Fatalf("sql.Open(): %s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("db.Ping(): %s", err)
	}
}
