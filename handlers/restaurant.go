package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Restaurant holds restaurant inforamtion
type Restaurant struct {
	RestID      int
	RestName    string
	RestAddress string
	RestLink    string
	RestOwner   string
}

// GetRestaurants returns restaurants list
func GetRestaurants() ([]*Restaurant, error) {
	rows, err := db.Query("SELECT RestID, RestName, RestAddress, RestLink, RestOwner FROM restaurants")
	if err != nil {
		log.Fatalf("db.Query(): %s", err)
		return nil, err
	}
	defer rows.Close()

	rests := make([]*Restaurant, 0)
	for rows.Next() {
		rest := new(Restaurant)
		err := rows.Scan(&rest.RestID, &rest.RestName, &rest.RestAddress, &rest.RestLink, &rest.RestOwner)
		if err != nil {
			log.Fatalf("rows.Scan(): %s", err)
			return nil, err
		}
		rests = append(rests, rest)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("rows.Err(): %s", err)
		return nil, err
	}
	return rests, nil
}

// GetRestName returns restaurant name
func GetRestName(request *http.Request) string {
	var restName string

	vars := mux.Vars(request)
	query := fmt.Sprintf("SELECT restName FROM restaurants WHERE restLink='%s'", vars["RestLink"])

	row := db.QueryRow(query)
	err := row.Scan(&restName)
	if err != nil {
		log.Fatal(err.Error())
		return ""
	}

	return restName
}
