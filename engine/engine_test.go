package engine

import (
	"os"
	"testing"
)

func TestCreateTheToDoListFileIfNeeded(t *testing.T) {
	toDoListFileName = "ToDoList_test_create.txt"
	defer os.Remove(toDoListFileName)

	os.Remove(toDoListFileName)
	created, err := createTheToDoListFileIfNeeded()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !created {
		t.Fatalf("Expected file to be created, but it wasn't")
	}

	file, err := os.Open(toDoListFileName)
	if err != nil {
		t.Fatalf("Expected file to exist, but got error: %v", err)
	}
	file.Close()
}

func TestGenerateItemId(t *testing.T) {
	toDoListFileName = "ToDoList_test_generate_id.txt"
	defer os.Remove(toDoListFileName)

	_, err := createTheToDoListFileIfNeeded()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	id, err := generateItemId(false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != 1 {
		t.Fatalf("Expected ID to be 1, got %d", id)
	}
}

func TestWriteAndReadExistingList(t *testing.T) {
	toDoListFileName = "ToDoList_test_write_read.txt"
	defer os.Remove(toDoListFileName)

	_, err := createTheToDoListFileIfNeeded()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	item := toDoItem{id: 1, name: "Test", description: "Test description"}
	err = writeItemToFile(item)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	list, err := readExistingList()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("Expected list length to be 1, got %d", len(list))
	}

	if list[0].name != "Test" || list[0].description != "Test description" {
		t.Fatalf("Expected item to match, got %+v", list[0])
	}
}
