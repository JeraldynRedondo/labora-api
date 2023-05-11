/*
Aquí definirás las estructuras de datos que representan tus objetos y también cualquier función
relacionada con la interacción con la base de datos.
*/
package model

import (
	"time"
)

var items []Item

type Item struct {
	Id            int       `json:"id"`
	Customer_name string    `json:"customerName"`
	Order_date    time.Time `json:"orderDate"`
	Product       string    `json:"product"`
	Quantity      int       `json:"quantity"`
	Price         int       `json:"price"`
}
