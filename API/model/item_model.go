package model //

import (
	"my_api_project/model"
	"time"
)

// Item is a struct that represents the Item object that belongs to the items table.
type Item struct {
	ID            int       `json:"id"`
	Customer_name string    `json:"customerName"`
	Order_date    time.Time `json:"orderDate"`
	Product       string    `json:"product"`
	Quantity      int       `json:"quantity"`
	Price         int       `json:"price"`
	Details       string    `json:"details"`
	TotalPrice    int       `json:"totalPrice"`
	ViewCount     int       `json:"viewCount"`
}

// CalculatedTotalPrice it is a function that returns the total price of an item.
func (item Item) CalculatedTotalPrice() int {
	totalPrice := item.Price * item.Quantity
	return totalPrice
}

// DBHandler is an interface that implements the methods of the database.
type DBHandler interface {
	GetItems() ([]model.Item, error)
	GetItemsPerPage(pages, itemsPerPage int) ([]model.Item, int, error)
	GetItemId(id int) (model.Item, error)
	GetItemName(name string) ([]model.Item, error)
	CreateItem(newItem model.Item)
	UpdateItem(id int, item model.Item) (model.Item, error)
	DeleteItem(id int) (model.Item, error)
	UpdateItemDetails(id int) (model.Item, error)
}
