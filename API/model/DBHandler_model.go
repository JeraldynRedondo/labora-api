package model

// DBHandler is an interface that implements the methods of the database.
type DBHandler interface {
	GetItems() ([]Item, error)
	GetItemsPerPage(pages, itemsPerPage int) ([]Item, int, error)
	GetItemId(id int) (Item, error)
	GetItemName(name string) ([]Item, error)
	CreateItem(newItem Item) (Item, error)
	UpdateItem(id int, item Item) (Item, error)
	DeleteItem(id int) error
}
