package interfaces

import (
	"todoapp/models"
)

type ToDoEngine interface {
	GetItems() ([]models.ToDoItem, error)
	CreateItem(name string, description string) error
	UpdateItem(id int, description string) error
	DeleteItem(id int) error
}
