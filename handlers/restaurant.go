package handlers

import "log"

// Restaurant holds restaurant inforamtion
type Restaurant struct {
	RestID      int
	RestName    string
	RestAddress string
	RestLink    string
}

// GetRestaurants returns restaurants list
func GetRestaurants() ([]*Restaurant, error) {
	rows, err := db.Query("SELECT RestID, RestName, RestAddress, RestLink FROM restaurants")
	if err != nil {
		log.Fatalf("db.Query(): %s", err)
		return nil, err
	}
	defer rows.Close()

	rests := make([]*Restaurant, 0)
	for rows.Next() {
		rest := new(Restaurant)
		err := rows.Scan(&rest.RestID, &rest.RestName, &rest.RestAddress, &rest.RestLink)
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
