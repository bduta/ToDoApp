// package main

// import (
// 	engine "todoapp/engine"

// 	"context"
// 	"log/slog"
// 	"os"

// 	"github.com/google/uuid"
// )

// func main() {
// 	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
// 	traceID := uuid.New().String()
// 	type contextKey string
// 	key := "TraceID"

// 	ctx := context.WithValue(context.Background(), contextKey(key), traceID)

// 	args := os.Args[1:]

// 	e := engine.NewEngine()
// 	err := e.ExecuteCommand(args)
// 	if err != nil {
// 		logger.With(key, ctx.Value(key)).Error(err.Error())
// 		return
// 	}
// }
