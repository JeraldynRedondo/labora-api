/*
Esta capa contendrá la lógica de negocio y se comunicará con la capa de modelos (en este caso,
con la base de datos) para recuperar, crear, actualizar y eliminar items.
*/
package service

import (
	"fmt"
	"my_api_project/model"
)

func GetItems() ([]model.Item, error) {
	items := make([]model.Item, 0)
	rows, err := Db.Query("SELECT * FROM items")
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.Id, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price)
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
		err := rows.Scan(&item.ID, &item.CustomerName, &item.OrderDate, &item.Product, &item.Quantity, &item.Price)
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

func createItem(newItem model.Item) (model.Item, error) {

	// Insertar el nuevo item en la base de datos
	insertStatement := `INSERT INTO items (customer_name, order_date, product, quantity, price)
                        VALUES ($1, $2, $3, $4, $5)`
	_, err := Db.Exec(insertStatement, newItem.Customer_name, newItem.Order_date, newItem.Product, newItem.Quantity, newItem.Price)
	if err != nil {
		return nil, err
	}
	return newItem, nil
}

func UpdateItem(id int, item model.Item) (model.Item, error) {
	var updatedItem model.Item
	row := Db.QueryRow("UPDATE items SET customer_name = $1, order_date = $2, product = $3, quantity = $4, price = $5 WHERE id = $6 RETURNING *",
		item.CustomerName, item.OrderDate, item.Product, item.Quantity, item.Price, id)
	err := row.Scan(&updatedItem.ID, &updatedItem.CustomerName, &updatedItem.OrderDate, &updatedItem.Product, &updatedItem.Quantity, &updatedItem.Price)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return updatedItem, nil
}

func DeleteItem(id int, item model.Item) (model.Item, error) {
	var deleteItem model.Item
	query := fmt.Sprintf("DELETE FROM items WHERE id = %d", id)
	_, err := Db.Exec(query)
	if err != nil {
		return nil, err
	}
	return deleteItem, nil
}
