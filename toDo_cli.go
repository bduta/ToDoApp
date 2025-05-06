package main

import (
	engine "engine"

	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	traceID := uuid.New().String()
	type contextKey string
	key := "TraceID"

	ctx := context.WithValue(context.Background(), contextKey(key), traceID)

	args := os.Args[1:]

	err := engine.ExecuteCommand(args)
	if err != nil {
		logger.With(key, ctx.Value(key)).Error(err.Error())
		return
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
