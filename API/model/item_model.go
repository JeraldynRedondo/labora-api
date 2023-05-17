package model //

import (
	"time"
)

var items []ItemDB

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

type ItemDB struct {
	Customer_name string    `json:"customerName"`
	Order_date    time.Time `json:"orderDate"`
	Product       string    `json:"product"`
	Quantity      int       `json:"quantity"`
	Price         int       `json:"price"`
	Details       string    `json:"details"`
}

func (item Item) CalculatedTotalPrice() int {
	totalPrice := item.Price * item.Quantity
	return totalPrice
}
