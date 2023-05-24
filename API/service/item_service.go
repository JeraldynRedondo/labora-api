package service

import "my_api_project/model"

type ItemService struct {
	DbHandler model.DBHandler
}

// GetItems implements the function GetItems of the database in a DBHandler.
func (s *ItemService) GetItems() ([]model.Item, error) {
	return s.DbHandler.GetItems()
}

// GetItemsPerPage implements the function GetItemsPerPage of the database in a DBHandler.
func (s *ItemService) GetItemsPerPage(pages, itemsPerPage int) ([]model.Item, int, error) {
	return s.DbHandler.GetItemsPerPage(pages, itemsPerPage)
}

// GetItemId implements the function GetItemId of the database in a DBHandler.
func (s *ItemService) GetItemId(id int) (model.Item, error) {
	return s.DbHandler.GetItemId(id)
}

// GetItemName implements the function GetItemName of the database in a DBHandler.
func (s *ItemService) GetItemName(name string) ([]model.Item, error) {
	return s.DbHandler.GetItemName(name)
}

// CreateItem implements the function CreateItem of the database in a DBHandler.
func (s *ItemService) CreateItem(newItem model.Item) (model.Item, error) {
	return s.DbHandler.CreateItem(newItem)
}

// UpdateItem implements the function UpdateItem of the database in a DBHandler.
func (s *ItemService) UpdateItem(id int, item model.Item) (model.Item, error) {
	return s.DbHandler.UpdateItem(id, item)
}

// DeleteItem implements the function DeleteItem of the database in a DBHandler.
func (s *ItemService) DeleteItem(id int) error {
	return s.DbHandler.DeleteItem(id)
}

/*
func (s *ItemService) UpdateItemDetails(id int) (model.Item, error) {
	return s.DbHandler.UpdateItemDetails(id)
}*/
