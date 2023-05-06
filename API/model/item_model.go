/*
Aquí definirás las estructuras de datos que representan tus objetos y también cualquier función
relacionada con la interacción con la base de datos.
*/
package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"my_api_project/config"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var items []Item

type Item struct {
	Id            int
	Customer_name string
	Order_date    time.Time
	Product       string
	Quantity      int
	Price         int
}

func raiz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Estás en la ruta raiz")
}

func getItems(w http.ResponseWriter, r *http.Request) {
	// Establecemos el encabezado "Content-Type" de la respuesta HTTP como "application/json"
	w.Header().Set("Content-Type", "application/json")

	dbConn, err := config.Connect_BD()
	query := `
		select id,	customer_name, order_date, product, quantity, price from items`
	rows, err := dbConn.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Id, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, item)
	}
	// Función para obtener todos los elementos
	json.NewEncoder(w).Encode(items)
}

func getItemId(w http.ResponseWriter, r *http.Request) {
	// Establecemos el encabezado "Content-Type" de la respuesta HTTP como "application/json"
	w.Header().Set("Content-Type", "application/json")

	dbConn, err := config.Connect_BD()

	// Función para obtener un elemento específico
	parametros := mux.Vars(r)
	id, err := strconv.Atoi(parametros["id"])
	if err != nil {
		fmt.Fprintf(w, "Inserte un item válido")
		return
	}

	var item Item
	err = dbConn.QueryRow("select id,	customer_name, order_date, product, quantity, price from items where id=$1", id).Scan(&item.Id, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price)
	if err != nil {
		fmt.Fprintf(w, "Error getting item from database")
		return
	}
	json.NewEncoder(w).Encode(item)

}

func getItemName(w http.ResponseWriter, r *http.Request) {
	// Establecemos el encabezado "Content-Type" de la respuesta HTTP como "application/json"
	w.Header().Set("Content-Type", "application/json")
	dbConn, err := config.Connect_BD()

	// Función para obtener un elemento específico
	parametros := mux.Vars(r)
	name := parametros["name"]

	var item Item
	err = dbConn.QueryRow("select id,	customer_name, order_date, product, quantity, price from items where customer_name=$1", name).Scan(&item.Id, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price)
	if err != nil {
		fmt.Fprintf(w, "Error getting item from database")
		return
	}
	json.NewEncoder(w).Encode(item)

}

func createItem(w http.ResponseWriter, r *http.Request) {
	// Función para crear un nuevo elemento
	var newItem Item
	rqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserte un item válido")
		return
	}

	json.Unmarshal(rqBody, &newItem)

	fmt.Println(newItem)

	dbConn, err := config.Connect_BD()

	// Insertar el nuevo item en la base de datos
	insertStatement := `INSERT INTO items (customer_name, order_date, product, quantity, price)
                        VALUES ($1, $2, $3, $4, $5)`
	_, err = dbConn.Exec(insertStatement, newItem.Customer_name, newItem.Order_date, newItem.Product, newItem.Quantity, newItem.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enviar una respuesta exitosa al cliente
	// Establecemos el encabezado "Content-Type" de la respuesta HTTP como "application/json"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
	fmt.Println("¡Item creado exitosamente!")
}
