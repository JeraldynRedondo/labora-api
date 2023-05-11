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
	"strings"
	"time"

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
	Items, err := service.GetItems()

	parameters := mux.Vars(request)
	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		// Manejar el error de la conversión
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("ID must be a number"))
		return
	}

	for _, item := range Items {
		if item.ID == id {
			json.NewEncoder(response).Encode(item)
			return
		}
	}
	json.NewEncoder(response).Encode(&model.Item{})
}

func GetItemByName(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	Items, err := service.GetItems()
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	parameters := mux.Vars(request)

	for _, item := range Items {
		if strings.EqualFold(item.Name, parameters["name"]) {
			json.NewEncoder(response).Encode(item)
			return
		}
	}
	json.NewEncoder(response).Encode(&model.Item{})
}

func CreateItem(response http.ResponseWriter, request *http.Request) {
	var newItem model.Item

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

	items := service.getItems()

	Json(response, http.StatusOK, items)
}

func UpdateItem(response http.ResponseWriter, request *http.Request) {
	Items, err := service.GetItems()
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

	for i, item := range Items {
		if item.ID == id {
			itemUpdate, err = service.UpdateItem(id, itemUpdate)
			if err != nil {
				http.Error(response, err.Error(), http.StatusBadRequest)
				return
			}
			items := service.getItems()
			Json(response, http.StatusOK, items)
			return
		}
	}

	response.Write([]byte("Could not update item"))
}

func DeleteItem(response http.ResponseWriter, request *http.Request) {
	Items, err := service.GetItems()
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

	for i, item := range Items {
		if item.ID == id {
			_, err = service.DeleteItem(id, item)
			if err != nil {
				http.Error(response, err.Error(), http.StatusBadRequest)
				return
			}
			items := service.getItems()
			Json(response, http.StatusOK, items)
			return
		}
	}
	response.Write([]byte("Item could not be removed"))
}

func EditarItemHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	idStr := vars["id"]

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("El ID debe ser un número"))
		return
	}

	// Leer el item a actualizar del body de la solicitud
	var itemToUpdate model.Item
	err = json.NewDecoder(request.Body).Decode(&itemToUpdate)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Error al procesar la solicitud"))
		return
	}

	// Parsear la fecha en el formato deseado
	t, err := time.Parse(time.RFC3339, itemToUpdate.OrderDate)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("La fecha debe estar en formato ISO 8601"))
		return
	}
	itemToUpdate.OrderDate = t.Format("2006-01-02")

	// Actualizar el item en la base de datos
	updatedItem, err := service.UpdateItemByID(id, itemToUpdate)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error al actualizar el item"))
		return
	}

	// Escribir el item actualizado en la respuesta
	Json(response, http.StatusOK, updatedItem)
}

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("You are on the root path")
}
