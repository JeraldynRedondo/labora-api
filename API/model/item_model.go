package model //

import (
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
