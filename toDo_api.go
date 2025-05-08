package main

import (
	"log"
	"net/http"

	engine "todoapp/engine"
	todoserver "todoapp/todoserver"
)

func main() {
	server := todoserver.NewToDoServer(engine.NewEngine())
	log.Fatal(http.ListenAndServe(":5000", server))
}
