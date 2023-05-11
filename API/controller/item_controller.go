/*
Aquí se ejecutan peticiones HTTP, como GET, POST, PUT y DELETE, y llamarás a los servicios
correspondientes.
*/

package controller

import (
	"encoding/json"
	"fmt"
	"math"
	"my_api_project/model"
	"my_api_project/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Json(response http.ResponseWriter, status int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("error while marshalling object %v, trace: %+v", data, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_, err = response.Write(bytes)
	if err != nil {
		fmt.Errorf("error while writing bytes to response writer: %+v", err)
	}
}

func GetAllItems(response http.ResponseWriter, _ *http.Request) {
	items, err := service.GetItems()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error getting items"))
		return
	}

	Json(response, http.StatusOK, items)
}

func GetItemsPaginated(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pageUser := r.URL.Query().Get("page")
	itemsUser := r.URL.Query().Get("itemsPerPage")

	page, err := strconv.Atoi(pageUser)
	if err != nil || page < 1 {
		page = 1
	}
	itemsPerPage, err := strconv.Atoi(itemsUser)
	if err != nil || itemsPerPage < 1 {
		itemsPerPage = 5
	}

	// Obtener la lista de elementos paginada
	newList, count, err := service.GetItemsPerPage(page, itemsPerPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	totalPages := int(math.Ceil(float64(count) / float64(itemsPerPage)))

	// Crear un mapa que contiene información sobre la paginación
	paginationInfo := map[string]interface{}{
		"totalPages":  totalPages,
		"currentPage": page,
	}

	// Crear un mapa que contiene la lista de elementos y la información de paginación
	responseData := map[string]interface{}{
		"items":      newList,
		"pagination": paginationInfo,
	}

	// Codificar el mapa de respuesta en formato JSON y enviar en la respuesta HTTP
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func GetItemById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var item model.Item
	parameters := mux.Vars(request)
	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		// Manejar el error de la conversión
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("ID must be a number"))
		return
	}

	item, err = service.GetItemId(id)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(response).Encode(item)
}

func GetItemByName(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	parameters := mux.Vars(request)
	name := parameters["name"]

	items, err := service.GetItemName(name)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(response).Encode(items)
}

func CreateItem(response http.ResponseWriter, request *http.Request) {
	var newItem model.Item
	var items []model.Item

	err := json.NewDecoder(request.Body).Decode(&newItem)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Error al procesar la solicitud"))
		return
	}

	newItem, err = service.CreateItem(newItem)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Error al procesar la solicitud"))
		return
	}

	items, err = service.GetItems()
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Error al procesar la solicitud"))
		return
	}

	Json(response, http.StatusOK, items)
}

func UpdateItem(response http.ResponseWriter, request *http.Request) {
	items, err := service.GetItems()
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	parameters := mux.Vars(request)
	var itemUpdate model.Item

	err = json.NewDecoder(request.Body).Decode(&itemUpdate)
	defer request.Body.Close()
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	itemUpdate, err = service.UpdateItem(id, itemUpdate)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	items, err = service.GetItems()
	Json(response, http.StatusOK, items)
}

func DeleteItem(response http.ResponseWriter, request *http.Request) {
	items, err := service.GetItems()
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	parameters := mux.Vars(request)
	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = service.DeleteItem(id)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	items, err = service.GetItems()
	Json(response, http.StatusOK, items)

}

func ItemDetails(response http.ResponseWriter, request *http.Request) {

	var updateItem model.Item
	parameters := mux.Vars(request)
	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		fmt.Println(err)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	updateItem, err = service.UpdateItemDetails(id)
	if err != nil {
		fmt.Println(err)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	Json(response, http.StatusOK, updateItem)

}

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("You are on the root path")
}
