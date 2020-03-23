package models

// Restaurants holds restaurants inforamtion
type Restaurants struct {
	restName            string
	restAddress         string
	restPictureLocation string
	restItems           Items
}

// Items holds restaurant items
type Items struct {
	itemName            string
	itemPrice           float32
	itemPictureLocation string
}

// GetMenu returns menu
func GetMenu() {

}
