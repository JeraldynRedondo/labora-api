package main

import (
	"log"
	"my_api_project/config"
	"my_api_project/controller"
	"my_api_project/service"

	"github.com/gorilla/mux"
)

func main() {

	service.Connect_DB()
	router := mux.NewRouter()

	router.HandleFunc("/", controller.Root).Methods("GET")
	router.HandleFunc("/items", controller.GetAllItems).Methods("GET")
	router.HandleFunc("/items/page", controller.GetItemsPaginated).Methods("GET")
	router.HandleFunc("/items/details/{id}", controller.ItemDetails).Methods("GET")
	router.HandleFunc("/items/id/{id}", controller.GetItemById).Methods("GET")
	router.HandleFunc("/items/name/{name}", controller.GetItemByName).Methods("GET")

	router.HandleFunc("/items", controller.CreateItem).Methods("POST")
	router.HandleFunc("/items/{id}", controller.UpdateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", controller.DeleteItem).Methods("DELETE")

	service.Db.PingOrDie()
	portNumber := ":3000"
	if err := config.StartServer(portNumber, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
