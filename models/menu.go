package models

// Restaurant holds restaurant inforamtion
type Restaurant struct {
	RestID      int
	RestName    string
	RestAddress string
	RestLink    string
	RestItems   []Item
}

// Item holds restaurant item information
type Item struct {
	ItemID    int
	ItemName  string
	ItemPrice float32
	ItemImage string
}

// GetRestaurant returns restaurants list
func GetRestaurant() ([]*Restaurant, error) {
	rows, err := db.Query("SELECT RestID, RestName, RestAddress, RestLink FROM restaurants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rests := make([]*Restaurant, 0)
	for rows.Next() {
		rest := new(Restaurant)
		err := rows.Scan(&rest.RestID, &rest.RestName, &rest.RestAddress, &rest.RestLink)
		if err != nil {
			return nil, err
		}
		rests = append(rests, rest)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rests, nil
}
