package models

import (
    "fmt"
    "log"
)

// Item holds restaurant item information
type Item struct {
	ItemID    int
	ItemName  string
	ItemPrice string
	ItemLink  string
    RestLink  string
}

// GetMenu returns menu list
func GetMenu(RestLink string) ([]*Item, error) {
    query := fmt.Sprintf("SELECT itemID, itemName, itemPrice, itemLink, restLink FROM items WHERE restLink='%s'", RestLink)
    rows, err := db.Query(query)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    defer rows.Close()

    items := make([]*Item, 0)
	for rows.Next() {
		item := new(Item)
		err := rows.Scan(&item.ItemID, &item.ItemName, &item.ItemPrice, &item.ItemLink, &item.RestLink)
		if err != nil {
            log.Fatal(err)
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
        log.Fatal(err)
		return nil, err
	}
	return items, nil
}
