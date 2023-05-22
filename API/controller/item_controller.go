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

// ResponseJson it is a function that sends the http response in Json format.
func ResponseJson(response http.ResponseWriter, status int, data interface{}) {
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

// GetItems it is a function that returns all the items from a request.
func GetAllItems(response http.ResponseWriter, _ *http.Request) {
	items, err := service.GetItems()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Error getting items, bad querying database"))
		return
	}

	ResponseJson(response, http.StatusOK, items)
}

// GetItemsPaginated it is a function that returns a number of items per page from a request.
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

	// Get the paginated list of items
	items, count, err := service.GetItemsPerPage(page, itemsPerPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	totalPages := int(math.Ceil(float64(count) / float64(itemsPerPage)))

	// Create a map containing information about pagination
	paginationInfo := map[string]interface{}{
		"totalPages":  totalPages,
		"currentPage": page,
	}

	// Create a map containing the list of items and the pagination information
	responseData := map[string]interface{}{
		"items":      items,
		"pagination": paginationInfo,
	}

	// Encode the response map in JSON format and send in the HTTP response
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

// GetItemById it is a function that returns the item that matches the id from a request.
func GetItemById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var item model.Item
	parameters := mux.Vars(request)
	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("ID must be a number"))
		return
	}

	item, err = service.GetItemId(id)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	ResponseJson(response, http.StatusOK, item)
}

// GetItemByName it is a function that returns the items that match the name from a request.
func GetItemByName(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	parameters := mux.Vars(request)
	name := parameters["name"]

	items, err := service.GetItemName(name)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	ResponseJson(response, http.StatusOK, items)
}

// CreateItem it is a function that creates an Item from a request.
func CreateItem(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var newItem model.Item

	err := json.NewDecoder(request.Body).Decode(&newItem)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Error processing the request"))
		return
	}

	newItem, err = service.CreateItem(newItem)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseJson(response, http.StatusOK, newItem)
}

// UpdateItem it is a function that updates an Item by id from a request.
func UpdateItem(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	parameters := mux.Vars(request)
	var itemUpdate model.Item

	err := json.NewDecoder(request.Body).Decode(&itemUpdate)
	defer request.Body.Close()
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("ID must be a number"))
		return
	}

	itemUpdate, err = service.UpdateItem(id, itemUpdate)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseJson(response, http.StatusOK, itemUpdate)
}

// DeleteItem it is a function that delete an Item by id from a request.
func DeleteItem(response http.ResponseWriter, request *http.Request) {
	items, err := service.GetItems()
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	parameters := mux.Vars(request)
	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("ID must be a number"))
		return
	}

	_, err = service.DeleteItem(id)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	items, err = service.GetItems()
	ResponseJson(response, http.StatusOK, items)

}

// ItemDetails it is a function that updates the Item details by id from a request.
func ItemDetails(response http.ResponseWriter, request *http.Request) {
	var updateItem model.Item

	parameters := mux.Vars(request)

	id, err := strconv.Atoi(parameters["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("ID must be a number"))
		return
	}

	updateItem, err = service.UpdateItemDetails(id)
	if err != nil {
		fmt.Println(err)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseJson(response, http.StatusOK, updateItem)
}

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("You are on the root path")
}
