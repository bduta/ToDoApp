package engine

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
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

	flagIdIndex := -1
	for index, item := range toDos {
		if item.id == flagId {
			flagIdIndex = index
		}
	}

	if flagIdIndex == -1 {
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
		return errors.New("could not open the ToDoList file")
	}
	defer f.Close()

	for _, item := range items {
		if _, err := f.WriteString(item.toFileFormat()); err != nil {
			return errors.New("could not append the toDo item to the ToDoList file: " + item.toFileFormat())
		}
	}

	return nil
}

func ExecuteCommand(arguments []string) error {
	if len(arguments) == 0 {
		toDos, err := readExistingList()
		if err != nil {
			return errors.New("Error reading existing list: " + err.Error())
		}

		for _, item := range toDos {
			fmt.Printf("Id:%d, ToDo:%s, Description:%s\n", item.id, item.name, item.description)
		}

		return nil
	}

	flag := arguments[0]
	switch strings.ToLower(flag) {
	case "-a":
		if len(arguments) != 3 {
			fmt.Println("Incorrect number of arguments. The format to add a flag is: -a FlagName FlagDescription")
			break
		}

		fileCreationRequired, err := createTheToDoListFileIfNeeded()
		if err != nil {
			return errors.New("Error creating ToDo list file: " + err.Error())
		}

		newItem := toDoItem{
			name:        arguments[1],
			description: arguments[2],
		}

		id, err := generateItemId(fileCreationRequired)
		if err != nil {
			return errors.New("Error generating item ID: " + err.Error())
		}
		newItem.id = id

		writeItemError := writeItemToFile(newItem)
		if writeItemError != nil {
			return errors.New("Error writing item to file: " + writeItemError.Error())
		}
	case "-u":
		if len(arguments) != 3 {
			fmt.Println("Incorrect number of arguments. The format to update a flag is: -u FlagId FlagDescription")
			break
		}

		toDos, err := readExistingList()
		if err != nil {
			return errors.New("Error reading existing list: " + err.Error())
		}

		index, err := getIndexBasedOnId(toDos, arguments[1])
		if err != nil {
			return errors.New("Error finding item by ID: " + err.Error())
		}

		toDos[index].description = arguments[2]

		overwritingFileErr := writeItemsToFile(toDos)
		if overwritingFileErr != nil {
			return errors.New("Error overwriting file: " + overwritingFileErr.Error())
		}

	case "-d":
		if len(arguments) != 2 {
			fmt.Println("Incorrect number of arguments. The format to delete a flag is: -d FlagId")
			break
		}

		toDos, err := readExistingList()
		if err != nil {
			return errors.New("Error reading existing list: " + err.Error())
		}

		index, err := getIndexBasedOnId(toDos, arguments[1])
		if err != nil {
			return errors.New("Error finding item by ID: " + err.Error())
		}

		toDos = append(toDos[:index], toDos[index+1:]...)

		overwritingFileErr := writeItemsToFile(toDos)
		if overwritingFileErr != nil {
			return errors.New("Error overwriting file: " + overwritingFileErr.Error())
		}
	default:
		fmt.Println("The flag entered is not valid.")
		fmt.Println("To add a flag: -a FlagName FlagDescription")
		fmt.Println("To update a flag: -u FlagId FlagDescription")
		fmt.Println("To delete a flag: -d FlagId")
	}
	return nil
}
