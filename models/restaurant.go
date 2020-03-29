package models

import "log"

// Restaurant holds restaurant inforamtion
type Restaurant struct {
	RestID      int
	RestName    string
	RestAddress string
	RestLink    string
}

// GetRestaurant returns restaurants list
func GetRestaurant() ([]*Restaurant, error) {
	rows, err := db.Query("SELECT RestID, RestName, RestAddress, RestLink FROM restaurants")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	rests := make([]*Restaurant, 0)
	for rows.Next() {
		rest := new(Restaurant)
		err := rows.Scan(&rest.RestID, &rest.RestName, &rest.RestAddress, &rest.RestLink)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		rests = append(rests, rest)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return rests, nil
}
