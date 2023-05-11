/*
Esta capa contendrá la lógica de negocio y se comunicará con la capa de modelos (en este caso,
con la base de datos) para recuperar, crear, actualizar y eliminar items.
*/
package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"my_api_project/model"
	"net/http"
	"strings"
)

func GetItems() ([]model.Item, error) {
	items := make([]model.Item, 0)
	rows, err := Db.Query("SELECT * FROM items")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func GetItemsPerPage(pages, itemsPerPage int) ([]model.Item, int, error) {
	// Calcular el índice inicial y el límite de elementos en función de la página actual y los elementos por página
	start := (pages - 1) * itemsPerPage

	// Obtener el número total de filas en la tabla items
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	// Obtener la lista de elementos correspondientes a la página actual
	rows, err := Db.Query("SELECT * FROM items ORDER BY id OFFSET $1 LIMIT $2", start, itemsPerPage)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var newListItems []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details)
		if err != nil {
			return nil, 0, err
		}
		newListItems = append(newListItems, item)
	}

	if len(newListItems) == 0 {
		return nil, 0, fmt.Errorf("No items found for page %d", pages)
	}
	return newListItems, count, nil
}

func GetItemId(id int) (model.Item, error) {
	var item model.Item
	err := Db.QueryRow("SELECT * FROM items WHERE id=$1", id).Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details)
	if err != nil {
		return model.Item{}, err
	}
	return item, nil
}

func GetItemName(name string) ([]model.Item, error) {
	var items []model.Item
	var item model.Item
	query := fmt.Sprintf("SELECT * FROM items WHERE customer_name LIKE '%%%s%%'", name)
	rows, err := Db.Query(query)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func CreateItem(newItem model.Item) (model.Item, error) {

	// Insertar el nuevo item en la base de datos
	insertStatement := `INSERT INTO items (customer_name, order_date, product, quantity, price)
                        VALUES ($1, $2, $3, $4, $5)`
	_, err := Db.Exec(insertStatement, newItem.Customer_name, newItem.Order_date, newItem.Product, newItem.Quantity, newItem.Price)
	if err != nil {
		return model.Item{}, err
	}
	return newItem, nil
}

func UpdateItem(id int, item model.Item) (model.Item, error) {
	var updatedItem model.Item
	row := Db.QueryRow("UPDATE items SET customer_name = $1, order_date = $2, product = $3, quantity = $4, price = $5 WHERE id = $6 RETURNING *",
		item.Customer_name, item.Order_date, item.Product, item.Quantity, item.Price, id)
	err := row.Scan(&updatedItem.ID, &updatedItem.Customer_name, &updatedItem.Order_date, &updatedItem.Product, &updatedItem.Quantity, &updatedItem.Price, &item.Details)
	if err != nil {
		return model.Item{}, err
	}
	return updatedItem, nil
}

func DeleteItem(id int) (model.Item, error) {
	var deleteItem model.Item
	query := fmt.Sprintf("DELETE FROM items WHERE id = %d", id)
	_, err := Db.Exec(query)
	if err != nil {
		return model.Item{}, err
	}
	return deleteItem, nil
}

func getWikipediaDetails(productName string) (string, error) {
	// Realizamos la petición a la API de Wikipedia
	url := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", productName)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	// Leemos la respuesta de la API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Decodificamos el JSON de la respuesta
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Obtenemos el extracto del objeto
	extract := data["extract"].(string)

	return extract, nil
}

func UpdateItemDetails(id int) (model.Item, error) {
	item, err := GetItemId(id)
	if err != nil {
		fmt.Println(err)
		return model.Item{}, err
	}
	product := strings.ToLower(item.Product)
	// Obtenemos el extracto del objeto desde la API de Wikipedia
	extract, err := getWikipediaDetails(product)
	if err != nil {
		fmt.Println(err)
		return model.Item{}, err
	}

	// Actualizamos la columna "detalles" en la tabla "items"
	var updatedItem model.Item
	row := Db.QueryRow("UPDATE items SET details=$1 WHERE id=$2 RETURNING *",
		extract, id)
	err = row.Scan(&updatedItem.ID, &updatedItem.Customer_name, &updatedItem.Order_date, &updatedItem.Product, &updatedItem.Quantity, &updatedItem.Price, &updatedItem.Details)
	if err != nil {
		fmt.Println(err)
		return model.Item{}, err
	}
	return updatedItem, nil
}
