package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var toDoListFileName string = "ToDoList.txt"

type toDoItem struct {
	id          int
	name        string
	description string
}

func (item toDoItem) toFileFormat() string {
	return fmt.Sprintf("%d,%s,%s\n", item.id, item.name, item.description)
}

func createTheToDoListFileIfNeeded() (bool, error) {
	creationRequired := false
	_, err := os.Stat(toDoListFileName)
	if os.IsNotExist(err) {
		_, err := os.Create(toDoListFileName)
		creationRequired = true
		if err != nil {
			return true, errors.New("file does not exist and could not be created")
		}
	}
	return creationRequired, nil
}

func readExistingList() (list []toDoItem, err error) {

	_, fileErr := os.Stat(toDoListFileName)
	if os.IsNotExist(fileErr) {
		return []toDoItem{}, errors.New("ToDo file does not exist")
	}

	file, fileErr := os.Open(toDoListFileName)
	if fileErr != nil {
		return []toDoItem{}, errors.New("file could not be opened")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var toDos []toDoItem
	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return []toDoItem{}, errors.New("Line has incorrect format: " + scanner.Text())
		}

		toDoId, err := strconv.Atoi(parts[0])
		if err != nil {
			return []toDoItem{}, errors.New("Id could not be converted to int: " + scanner.Text())
		}

		toDo := toDoItem{
			id:          toDoId,
			name:        parts[1],
			description: parts[2],
		}
		toDos = append(toDos, toDo)
	}

	slices.SortFunc(toDos, func(i, j toDoItem) int {
		return i.id - j.id
	})

	return toDos, nil
}

func generateItemId(fileCreationRequired bool) (id int, err error) {
	if !fileCreationRequired {
		toDos, err := readExistingList()
		if err != nil {
			return -1, err
		}

		if len(toDos) > 0 {
			lastToDo := toDos[len(toDos)-1]
			return lastToDo.id + 1, nil
		} else {
			return 1, nil
		}

	} else {
		return 1, nil
	}
}

func getIndexBasedOnId(toDos []toDoItem, id string) (index int, err error) {
	flagId, err := strconv.Atoi(id)
	if err != nil {
		return -1, errors.New("Id could not be converted to int: " + id)
	}

	flagIdIndex := sort.Search(len(toDos), func(i int) bool {
		return toDos[i].id == flagId
	})

	if flagIdIndex == len(toDos) {
		return -1, errors.New("flag Id could not be found")
	}

	return flagIdIndex, nil
}

func writeItemToFile(item toDoItem) error {

	_, err := os.Stat(toDoListFileName)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("ToDoList file does not exist")
	} else {

		f, err := os.OpenFile(toDoListFileName,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		if err != nil {
			return errors.New("could not open the ToDoList file")
		}

		if _, err := f.WriteString(item.toFileFormat()); err != nil {
			return errors.New("could not append the new toDo item to the ToDoList file: " + item.toFileFormat())
		}
	}

	return nil
}

func writeItemsToFile(items []toDoItem) error {
	f, err := os.Create(toDoListFileName)
	if err != nil {
		err = errors.New("could not open the ToDoList file")
	}
	defer f.Close()

	for _, item := range items {
		if _, err := f.WriteString(item.toFileFormat()); err != nil {
			return errors.New("could not append the toDo item to the ToDoList file: " + item.toFileFormat())
		}
	}

	return nil
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		toDos, err := readExistingList()
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, item := range toDos {
			fmt.Printf("Id:%d, ToDo:%s, Description:%s\n", item.id, item.name, item.description)
		}

		return
	}

	flag := args[0]
	switch strings.ToLower(flag) {
	case "-a":
		if len(args) != 3 {
			fmt.Println("Incorrect number of arguments. The format to add a flag is: -a FlagName FlagDescription")
			break
		}

		fileCreationRequired, err := createTheToDoListFileIfNeeded()
		if err != nil {
			log.Fatal(err)
			return
		}

		newItem := toDoItem{
			name:        args[1],
			description: args[2],
		}

		id, err := generateItemId(fileCreationRequired)
		if err != nil {
			log.Fatal(err)
			return
		}
		newItem.id = id

		writeItemError := writeItemToFile(newItem)
		if writeItemError != nil {
			log.Fatal(writeItemError)
			return
		}
	case "-u":
		if len(args) != 3 {
			fmt.Println("Incorrect number of arguments. The format to update a flag is: -u FlagId FlagDescription")
			break
		}

		toDos, err := readExistingList()
		if err != nil {
			log.Fatal(err)
			return
		}

		index, err := getIndexBasedOnId(toDos, args[1])
		if err != nil {
			log.Fatal(err)
			return
		}

		toDos[index].description = args[2]

		overwritingFileErr := writeItemsToFile(toDos)
		if overwritingFileErr != nil {
			log.Fatal(overwritingFileErr)
			return
		}

	case "-d":
		if len(args) != 2 {
			fmt.Println("Incorrect number of arguments. The format to delete a flag is: -d FlagId")
			break
		}

		toDos, err := readExistingList()
		if err != nil {
			log.Fatal(err)
			return
		}

		index, err := getIndexBasedOnId(toDos, args[1])
		if err != nil {
			log.Fatal(err)
			return
		}

		toDos = append(toDos[:index], toDos[index+1:]...)

		overwritingFileErr := writeItemsToFile(toDos)
		if overwritingFileErr != nil {
			log.Fatal(overwritingFileErr)
			return
		}
	default:
		fmt.Println("The flag entered is not valid.")
		fmt.Println("To add a flag: -a FlagName FlagDescription")
		fmt.Println("To update a flag: -u FlagId FlagDescription")
		fmt.Println("To delete a flag: -d FlagId")
	}
}

/*
	File does not exist, add item - tested
	File does not exist, update item with correct number of arguments - tested
	File does not exist, update item with incorrect number of arguments - tested
	File does not exist, delete item with correct number of arguments - tested
	File does not exist, delete item with incorrect number of arguments - tested
	File does not exist, no flag - tested
	File does not exist, invalid flag - tested


	File exists, add item - tested
	File exists, add item with incorrect number of arguments - tested

	File exists, no flag - tested

	File exists, update item with existing id at the top of the file
	File exists, update item with existing id at the middle of the file
	File exists, update item with existing id at the end of the file
	File exists, update item with non-existing id in the file
	File exists, update item with id which is not int
	File exists, update item with incorrect number of arguments

	File exists, delete item with existing id at the top of the file
	File exists, delete item with existing id at the middle of the file
	File exists, delete item with existing id at the end of the file
	File exists, delete item with non-existing id in the file
	File exists, delete item with id which is not int
	File exists, delete item with incorrect number of arguments

	File exists, invalid flag - tested
*/
