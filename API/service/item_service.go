package service

import (
	"fmt"
	"io/ioutil"
	"my_api_project/model"
	"net/http"
	"sync"
	"time"
)

func GetItems() ([]model.Item, error) {
	items := make([]model.Item, 0)
	rows, err := Db.Query("SELECT * FROM items ORDER BY id")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details, &item.TotalPrice, &item.ViewCount)
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
		err := rows.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details, &item.TotalPrice, &item.ViewCount)
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
	UpdateViewCount(id)
	err := Db.QueryRow("SELECT * FROM items WHERE id=$1", id).Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details, &item.TotalPrice, &item.ViewCount)
	if err != nil {
		return model.Item{}, err
	}
	return item, nil
}

func GetItemName(name string) ([]model.Item, error) {
	var items []model.Item
	var item model.Item
	query := fmt.Sprintf("SELECT * FROM items WHERE customer_name ILIKE '%%%s%%'", name)
	rows, err := Db.Query(query)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details, &item.TotalPrice, &item.ViewCount)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func CreateItem(newItem model.Item) (model.Item, error) {
	details, err := getDetails()
	if err != nil {
		return model.Item{}, err
	}
	newItem.TotalPrice = newItem.CalculatedTotalPrice()
	// Insertar el nuevo item en la base de datos
	insertStatement := `INSERT INTO items (customer_name, order_date, product, quantity, price,details,total_price ,view_count)
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	row := Db.QueryRow(insertStatement, newItem.Customer_name, time.Now(), newItem.Product, newItem.Quantity, newItem.Price, details, newItem.TotalPrice, 0)

	err = row.Scan(&newItem.ID, &newItem.Customer_name, &newItem.Order_date, &newItem.Product, &newItem.Quantity, &newItem.Price, &newItem.Details, &newItem.TotalPrice, &newItem.ViewCount)
	if err != nil {
		return model.Item{}, err
	}
	if err != nil {
		return model.Item{}, err
	}

	return newItem, nil
}

func UpdateItem(id int, item model.Item) (model.Item, error) {
	var updatedItem model.Item

	details, err := getDetails()
	if err != nil {
		return model.Item{}, err
	}
	item.TotalPrice = item.CalculatedTotalPrice()

	row := Db.QueryRow("UPDATE items SET customer_name = $1, order_date = $2, product = $3, quantity = $4, price = $5, details = $6, total_price=$7 WHERE id = $8 RETURNING *",
		item.Customer_name, time.Now(), item.Product, item.Quantity, item.Price, details, item.TotalPrice, id)
	err = row.Scan(&updatedItem.ID, &updatedItem.Customer_name, &updatedItem.Order_date, &updatedItem.Product, &updatedItem.Quantity, &updatedItem.Price, &updatedItem.Details, &updatedItem.TotalPrice, &updatedItem.ViewCount)
	if err != nil {
		return model.Item{}, err
	}
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

func getDetails() (string, error) {
	// Realizamos la petición a la API de loripsum
	url := fmt.Sprintf("http://loripsum.net/api/1/short")
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

	data := string(body)
	return data, nil
}

func UpdateItemDetails(id int) (model.Item, error) {
	// Obtenemos el párrafo del objeto desde la API de loripsum
	detail, err := getDetails()
	if err != nil {
		fmt.Println(err)
		return model.Item{}, err
	}

	// Actualizamos la columna "detalles" en la tabla "items"
	var updatedItem model.Item
	row := Db.QueryRow("UPDATE items SET details=$1 WHERE id=$2 RETURNING *",
		detail, id)
	err = row.Scan(&updatedItem.ID, &updatedItem.Customer_name, &updatedItem.Order_date, &updatedItem.Product, &updatedItem.Quantity, &updatedItem.Price, &updatedItem.Details, &updatedItem.TotalPrice, &updatedItem.ViewCount)
	if err != nil {
		fmt.Println(err)
		return model.Item{}, err
	}
	return updatedItem, nil
}

// Funcio que actualiza los precios totales de los items
func UpdateTotalPriceItem(item model.Item) (model.Item, error) {
	totalPrice := item.CalculatedTotalPrice()
	//fmt.Printf("El precio total del item %d con precio de %d y cantidad %d es %d \n", item.ID, item.Price, item.Quantity, totalPrice)
	row := Db.QueryRow("UPDATE items SET total_price=$1 WHERE id = $2 RETURNING *", totalPrice, item.ID)
	err := row.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details, &item.TotalPrice, &item.ViewCount)
	if err != nil {
		return model.Item{}, err
	}
	return item, nil
}

var m sync.Mutex

func UpdateViewCount(id int) {
	m.Lock()
	Db.QueryRow("UPDATE items SET view_count = view_count + $1 WHERE id = $2 RETURNING *",
		1, id)
	m.Unlock()
}
