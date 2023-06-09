package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func raiz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Estás en la ruta raiz")
}

func getItems(w http.ResponseWriter, r *http.Request) {
	// Establecemos el encabezado "Content-Type" de la respuesta HTTP como "application/json"
	w.Header().Set("Content-Type", "application/json")

	// Función para obtener todos los elementos
	json.NewEncoder(w).Encode(items)
}

func getItemId(w http.ResponseWriter, r *http.Request) {
	// Establecemos el encabezado "Content-Type" de la respuesta HTTP como "application/json"
	w.Header().Set("Content-Type", "application/json")

	// Función para obtener un elemento específico
	parametros := mux.Vars(r)

	for _, item := range items {
		if item.ID == parametros["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	fmt.Fprintf(w, "Usuario no encontrado")
	json.NewEncoder(w).Encode(&Item{})

}

func getItemName(w http.ResponseWriter, r *http.Request) {
	// Establecemos el encabezado "Content-Type" de la respuesta HTTP como "application/json"
	w.Header().Set("Content-Type", "application/json")

	// Función para obtener un elemento específico
	parametros := mux.Vars(r)
	var itemsIguales []Item
	count := 0
	for _, item := range items {
		if strings.EqualFold(item.Name, parametros["name"]) {
			itemsIguales = append(itemsIguales, item)
			count++
		}
	}

	if count > 1 {
		json.NewEncoder(w).Encode(itemsIguales)
		return
	}

	fmt.Fprintf(w, "Usuario no encontrado")
	json.NewEncoder(w).Encode(&Item{})

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

	id := len(items) + 1
	// Agregar el nuevo elemento al slice
	newItem.ID = "item" + strconv.Itoa(id)
	items = append(items, newItem)

	// Enviar una respuesta exitosa al cliente
	// Establecemos el encabezado "Content-Type" de la respuesta HTTP como "application/json"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
	fmt.Println("¡Item creadoexitosamente!")
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	// Función para actualizar un elemento existente
	var UpdateItem Item
	parametros := mux.Vars(r)
	rqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserte un item válido")
		return
	}
	json.Unmarshal(rqBody, &UpdateItem)

	for i, item := range items {
		if item.ID == parametros["id"] {
			items = append(items[:i], items[i+1:]...)
			UpdateItem.ID = parametros["id"]
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

	for i, item := range items {
		if item.ID == parametros["id"] {
			items = append(items[:i], items[i+1:]...)
			fmt.Fprintf(w, "La tarea con el ID %v fue eliminada", item.ID)
			return
		}
	}

	fmt.Fprintf(w, "Usuario no encontrado")
}

// var items = []Item{{ID:   "1",Name: "Juana",},{ID:   "2",Name: "Mario",},{ID:   "3",Name: "Paola",},{ID:   "4",Name: "Luis",},
// {ID:   "5",Name: "Isabella",},{ID:   "6",Name: "Jose",},{ID:   "7",Name: "Elena",},{ID:   "8",Name: "Pedro",},{ID:   "9",
// Name: "Laura",},{ID:   "10",Name: "Samuel",},}

var items []Item

func main() {

	for i := 1; i <= 10; i++ {
		items = append(items, Item{ID: fmt.Sprintf("item%d", i), Name: fmt.Sprintf("Item %d", i)})
	}

	//items = append(items, Item{ID: fmt.Sprintf("item%d", 11), Name: fmt.Sprintf("Item %d", 1)})

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
	err := http.ListenAndServe(":"+strconv.Itoa(portNumber), router)
	if err != nil {
		fmt.Println(err)
	}

}
