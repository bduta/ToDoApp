package interfaces

type ToDoEngine interface {
	GetItems() ([]byte, error)
	CreateItem(name string, description string) error
	UpdateItem(id int, description string) error
	DeleteItem(id int) error
}
