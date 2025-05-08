package todoserver

import (
	"encoding/json"
	"net/http"

	"todoapp/interfaces"
	"todoapp/models"
)

type ToDoServer struct {
	toDoEngine interfaces.ToDoEngine
	http.Handler
}

const jsonContentType = "application/json"

func NewToDoServer(toDoEngine interfaces.ToDoEngine) *ToDoServer {
	t := new(ToDoServer)

	t.toDoEngine = toDoEngine

	router := http.NewServeMux()
	router.Handle("/create", http.HandlerFunc(t.createHandler))
	router.Handle("/get", http.HandlerFunc(t.getHandler))
	router.Handle("/update", http.HandlerFunc(t.updateHandler))
	router.Handle("/delete", http.HandlerFunc(t.deleteHandler))

	t.Handler = router

	return t
}

// createHandler handles the creation of a new ToDo item.
func (t *ToDoServer) createHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed"}`))
		return
	}

	var input models.ToDoItem

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}

	error := t.toDoEngine.CreateItem(input.Name, input.Description)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to create ToDo item: ` + error.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "ToDo item created successfully"}`))
}

func (t *ToDoServer) getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed"}`))
		return
	}

	toDosJson, err := t.toDoEngine.GetItems()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to get ToDo items: ` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(toDosJson)
}

func (t *ToDoServer) updateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed"}`))
		return
	}

	var input models.ToDoItem

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}

	error := t.toDoEngine.UpdateItem(input.Id, input.Description)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to update ToDo item: ` + error.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ToDo item updated successfully"}`))
}

func (t *ToDoServer) deleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed"}`))
		return
	}

	var input models.ToDoItem

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}

	error := t.toDoEngine.DeleteItem(input.Id)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to delete ToDo item: ` + error.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ToDo item deleted successfully"}`))
}
