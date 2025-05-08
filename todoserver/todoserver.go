package todoserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

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
	router.Handle("/create", traceMiddleware(http.HandlerFunc(t.createHandler)))
	router.Handle("/get", traceMiddleware(http.HandlerFunc(t.getHandler)))
	router.Handle("/update", traceMiddleware(http.HandlerFunc(t.updateHandler)))
	router.Handle("/delete", traceMiddleware(http.HandlerFunc(t.deleteHandler)))

	t.Handler = router

	return t
}

func traceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), "traceID", traceID)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (t *ToDoServer) createHandler(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed", "traceID": "` + traceID + `"}`))
		return
	}

	var input models.ToDoItem

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload", "traceID": "` + traceID + `"}`))
		return
	}

	error := t.toDoEngine.CreateItem(input.Name, input.Description)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to create ToDo item: ` + error.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "ToDo item created successfully", "traceID": "` + traceID + `"}`))
}

func (t *ToDoServer) getHandler(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed", "traceID": "` + traceID + `"}`))
		return
	}

	toDosJson, err := t.toDoEngine.GetItems()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to get ToDo items: ` + err.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(toDosJson)
}

func (t *ToDoServer) updateHandler(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed", "traceID": "` + traceID + `"}`))
		return
	}

	var input models.ToDoItem

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload", "traceID": "` + traceID + `"}`))
		return
	}

	error := t.toDoEngine.UpdateItem(input.Id, input.Description)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to update ToDo item: ` + error.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ToDo item updated successfully", "traceID": "` + traceID + `"}`))
}

func (t *ToDoServer) deleteHandler(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed", "traceID": "` + traceID + `"}`))
		return
	}

	var input models.ToDoItem

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload", "traceID": "` + traceID + `"}`))
		return
	}

	error := t.toDoEngine.DeleteItem(input.Id)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to delete ToDo item: ` + error.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ToDo item deleted successfully", "traceID": "` + traceID + `"}`))
}
