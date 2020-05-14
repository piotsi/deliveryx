package handlers

import (
	"fmt"
	"log"
)

// Item holds restaurant item information
type Item struct {
	ItemID          int
	ItemName        string
	ItemPrice       string
	ItemLink        string
	ItemDescription string
	RestLink        string
}

// GetMenu returns menu list
func GetMenu(RestLink string) ([]*Item, error) {
	query := fmt.Sprintf("SELECT itemID, itemName, itemPrice, itemLink, itemDescription, restLink FROM items WHERE restLink='%s'", RestLink)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("db.Query(): %s", err)
		return nil, err
	}
	defer rows.Close()

	items := make([]*Item, 0)
	for rows.Next() {
		item := new(Item)
		err := rows.Scan(&item.ItemID, &item.ItemName, &item.ItemPrice, &item.ItemLink, &item.ItemDescription, &item.RestLink)
		if err != nil {
			log.Fatalf("rows.Scan(): %s", err)
			return nil, err
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("rows.Err(): %s", err)
		return nil, err
	}
	return items, nil
}
