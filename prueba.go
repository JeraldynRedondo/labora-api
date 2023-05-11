package api

/*
package main

import (
	_ "github.com/lib/pq"
)

/*
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", raiz).Methods("GET")
	router.HandleFunc("/items", getItems).Methods("GET")
	router.HandleFunc("/items/id/{id}", getItemId).Methods("GET")
	router.HandleFunc("/items/name/{name}", getItemName).Methods("GET")
	router.HandleFunc("/items", createItem).Methods("POST")
	router.HandleFunc("/items/{id}", updateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

	var portNumber int = 9999
	fmt.Println("Listen in port ", portNumber)
	errP := http.ListenAndServe(":"+strconv.Itoa(portNumber), router)
	if errP != nil {
		fmt.Println(errP)
	}

}

const (
	host        = "localhost"
	port        = "5432"
	dbName      = "labora-proyect-1"
	rolName     = "postgres"
	rolPassword = "1234"
)

func Connect_BD() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, rolName, rolPassword, dbName)
	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successful connection to the database:", dbConn)
	return dbConn, err
}

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

	dbConn, err := Connect_BD()
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

	dbConn, err := Connect_BD()

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
	dbConn, err := Connect_BD()

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

	dbConn, err := Connect_BD()

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

func updateItem(w http.ResponseWriter, r *http.Request) {
	// Función para actualizar un elemento existente
	var UpdateItem Item
	parametros := mux.Vars(r)
	id, err := strconv.Atoi(parametros["id"])
	if err != nil {
		fmt.Fprintf(w, "Inserte un item válido")
		return
	}
	rqBody, err2 := ioutil.ReadAll(r.Body)
	if err2 != nil {
		fmt.Fprintf(w, "Inserte un item válido")
		return
	}
	json.Unmarshal(rqBody, &UpdateItem)

	for i, item := range items {
		if item.Id == id {
			items = append(items[:i], items[i+1:]...)
			UpdateItem.Id = id
			items = append(items, UpdateItem)
		}
	}

	// Enviar una respuesta exitosa al cliente
	// Establecemos el encabezado "Content-Type" de la respuesta HTTP como "application/json"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UpdateItem)
	fmt.Println("¡Item actualizado!")
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	// Establecemos el encabezado "Content-Type" de la respuesta HTTP como "application/json"
	w.Header().Set("Content-Type", "application/json")

	// Función para eliminar un elemento
	parametros := mux.Vars(r)
	id, err := strconv.Atoi(parametros["id"])
	if err != nil {
		fmt.Fprintf(w, "Inserte un item válido")
		return
	}

	for i, item := range items {
		if item.Id == id {
			items = append(items[:i], items[i+1:]...)
			fmt.Fprintf(w, "La tarea con el ID %v fue eliminada", item.Id)
			return
		}
	}

	fmt.Fprintf(w, "Usuario no encontrado")
}
*/
